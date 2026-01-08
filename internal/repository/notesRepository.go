package repository

import (
	"blog-service/internal/models"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type NotesRepoImpl struct {
	db *sql.DB
}

func NewNotesRepoImpl(db *sql.DB) *NotesRepoImpl {
	return &NotesRepoImpl{
		db: db,
	}
}

func (repo *NotesRepoImpl) Save(note models.NoteDto, userId int) (int, error) {
	op := "internal.repository.NotesRepoImpl.SaveNote"
	var insertId int
	err := repo.db.QueryRow("INSERT INTO note(title, content, create_at, updated_at, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING note_id", note.Title, note.Content, note.CreatedAt, note.UpdatedAt, userId).Scan(&insertId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return insertId, nil
}

func (repo *NotesRepoImpl) UpdateNote(note models.NoteDto) error {
	op := "internal.repository.NotesRepoImpl.UpdateNote"
	err := repo.db.QueryRow("UPDATE note SET title = $1, content = $2, updated_at = $3 WHERE note_id = $4", note.Title, note.Content, time.Now(), note.NoteId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (repo *NotesRepoImpl) GetById(id int) (models.NoteDto, error) {
	op := "internal.repository.NotesRepoImpl.GetById"
	var note models.NoteDto
	err := repo.db.QueryRow("SELECT note.note_id, note.content, note.title, note.updated_at, note.create_at, note.user_id FROM note WHERE note_id = $1", id).Scan(&note.NoteId, &note.Content, &note.Title, &note.UpdatedAt, &note.CreatedAt, &note.UserId)
	if err != nil {
		fmt.Printf("%s: %s\n", op, err)
		return models.NoteDto{}, err
	}
	return note, nil
}

func (repo *NotesRepoImpl) DeleteById(id int) error {
	op := "internal.repository.NotesRepoImpl.DeleteById"
	_, err := repo.db.Exec("DELETE FROM note WHERE note_id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (repo *NotesRepoImpl) GetAllNotes(userId int, page models.PageDto) ([]models.NoteDto, error) {
	op := "internal.repository.NotesRepoImpl.GetAllNotes"
	var notes []models.NoteDto
	query := "SELECT note_id, title, content, create_at, updated_at, user_id FROM note WHERE user_id = $1 ORDER BY create_at " + page.Sort + " OFFSET $2 LIMIT $3"
	rows, err := repo.db.Query(query, userId, page.Offset, page.Limit)
	if err != nil {
		fmt.Printf("%s: %s\n", op, err)
		return notes, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("%s: %s\n", op, err)
		}
	}(rows)

	for rows.Next() {
		var note models.NoteDto
		if err := rows.Scan(&note.NoteId, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.UserId); err != nil {
			log.Printf("%s: %s\n", op, err)
			return nil, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%s: %s\n", op, err)
		return nil, err
	}

	return notes, nil
}
