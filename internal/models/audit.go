package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Audit struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"user_id" json:"user_id"`
	Action    string        `bson:"action" json:"action"`
	Changes   string        `bson:"changes" json:"changes"` // JSON string of changes
	Type      string        `bson:"type" json:"type"`
	TypeID    bson.ObjectID `bson:"type_id" json:"type_id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	_OldData  any           `bson:"-" json:"-"`
	_NewData  any           `bson:"-" json:"-"`
}

type FieldChange struct {
	Field    string `json:"field"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
	Action   string `json:"action"` // "added", "modified", "removed"
}

type AuditChanges struct {
	Summary      string        `json:"summary"`
	TotalChanges int           `json:"total_changes"`
	Fields       []FieldChange `json:"fields"`
}

func (a *Audit) BeforeAction(data any) {
	a._OldData = data
}

func (a *Audit) AfterAction(data any) {
	a._NewData = data

	// Compare old and new data to create audit entry
	changes := a.createDetailedChanges()

	// Convert changes to JSON string for easy reading
	changesJSON, _ := json.MarshalIndent(changes, "", "  ")
	a.Changes = string(changesJSON)
}

func (a *Audit) createDetailedChanges() AuditChanges {
	changes := AuditChanges{
		Fields: []FieldChange{},
	}

	// Convert both to maps for easier comparison
	oldMap := a.structToMap(a._OldData)
	newMap := a.structToMap(a._NewData)

	// Find all unique field names
	allFields := make(map[string]bool)
	for field := range oldMap {
		allFields[field] = true
	}
	for field := range newMap {
		allFields[field] = true
	}

	// Compare each field
	for field := range allFields {
		oldValue, oldExists := oldMap[field]
		newValue, newExists := newMap[field]

		if !oldExists && newExists {
			// Field was added
			changes.Fields = append(changes.Fields, FieldChange{
				Field:    field,
				OldValue: nil,
				NewValue: newValue,
				Action:   "added",
			})
		} else if oldExists && !newExists {
			// Field was removed
			changes.Fields = append(changes.Fields, FieldChange{
				Field:    field,
				OldValue: oldValue,
				NewValue: nil,
				Action:   "removed",
			})
		} else if oldExists && newExists {
			// Check if field was modified
			if !cmp.Equal(oldValue, newValue) {
				changes.Fields = append(changes.Fields, FieldChange{
					Field:    field,
					OldValue: oldValue,
					NewValue: newValue,
					Action:   "modified",
				})
			}
		}
	}

	changes.TotalChanges = len(changes.Fields)
	changes.Summary = a.generateSummary(changes.Fields)

	return changes
}

func (a *Audit) structToMap(data any) map[string]any {
	if data == nil {
		return make(map[string]any)
	}

	// Convert to JSON and back to map for easier handling
	jsonBytes, _ := json.Marshal(data)
	var result map[string]any
	json.Unmarshal(jsonBytes, &result)

	return result
}

func (a *Audit) generateSummary(fields []FieldChange) string {
	if len(fields) == 0 {
		return "No changes detected"
	}

	added := 0
	modified := 0
	removed := 0

	for _, field := range fields {
		switch field.Action {
		case "added":
			added++
		case "modified":
			modified++
		case "removed":
			removed++
		}
	}

	summary := ""
	if added > 0 {
		summary += fmt.Sprintf("%d field(s) added", added)
	}
	if modified > 0 {
		if summary != "" {
			summary += ", "
		}
		summary += fmt.Sprintf("%d field(s) modified", modified)
	}
	if removed > 0 {
		if summary != "" {
			summary += ", "
		}
		summary += fmt.Sprintf("%d field(s) removed", removed)
	}

	return summary
}
