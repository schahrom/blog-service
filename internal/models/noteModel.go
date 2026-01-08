package models

import "time"

type NoteDto struct {
	NoteId    int       `json:"note_id"`
	Title     string    `json:"title,omitempty" validate:"required"`
	Content   string    `json:"content,omitempty" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    int       `json:"user_id"`
}
