package _type

import "github.com/uchupx/saceri-chatbot-api/pkg/helper"

type GetQuery struct {
	Keyword *string `query:"keyword"`
	Page    *int    `query:"page"`
	PerPage *int    `query:"per_page"`
}

type UserUpdateRequest struct {
	id   string  `json:"-"`
	Name *string `json:"name"`
	//Password *string `json:"password,omitempty"`
}

func (u *UserUpdateRequest) SetID(id string) {
	u.id = id
}

func (u UserUpdateRequest) ID() string {
	return u.id
}

func (q GetQuery) Limit() int {
	return helper.DefaultInt(q.PerPage, 10)
}

func (q GetQuery) Offset() int {
	page := helper.DefaultInt(q.Page, 1)
	return (page - 1) * q.Limit()
}
