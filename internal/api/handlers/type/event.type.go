package _type

type EventCreateUpdateRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Place          string `json:"place"`
	DonationTarget *uint  `json:"donation_target"`
	StartAt        string `json:"start_at"`
	EndAt          string `json:"end_at"`
}
