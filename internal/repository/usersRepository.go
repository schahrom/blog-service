package repository

import (
	"blog-service/internal/lib/hasher"
	"blog-service/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type UsersRepositoryImpl struct {
	db     *sql.DB
	hasher hasher.PasswordHasher
}

func NewUsersRepositoryImpl(db *sql.DB, hasher hasher.PasswordHasher) *UsersRepositoryImpl {
	return &UsersRepositoryImpl{db: db, hasher: hasher}
}

func (r *UsersRepositoryImpl) Create(user *models.User) (models.User, error) {
	const op = "internal.repository.UserRepositoryImpl.Create"
	tx, err := r.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	var userId int
	hashedPassword := r.hasher.Hash(user.Password)
	err = r.db.QueryRow("INSERT INTO \"user\" (fist_name, last_name, login, email, password, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id",
		user.FirstName,
		user.LastName,
		user.Login,
		user.Email,
		hashedPassword,
		time.Now(),
	).Scan(&userId)
	if err != nil {
		fmt.Printf("%s: %s\n", op, err)
		return models.User{}, err
	}
	return r.GetById(userId)
}

func (r *UsersRepositoryImpl) GetById(id int) (models.User, error) {
	const op = "internal.repository.UserRepositoryImpl"
	var respUser models.User
	err := r.db.QueryRow(
		"SELECT * FROM \"user\" WHERE user_id = $1",
		id,
	).Scan(
		&respUser.UserId,
		&respUser.FirstName,
		&respUser.LastName,
		&respUser.Login,
		&respUser.Email,
		&respUser.Password,
		&respUser.CreatedAt,
	)
	if err != nil {
		fmt.Printf("%s: %s\n", op, err)
		return models.User{}, err
	}

	return respUser, nil
}

func (r *UsersRepositoryImpl) GetByCredentials(login string, hashedPassword string) (models.User, error) {
	const op = "internal.repository.UserRepositoryImpl"
	var respUser models.User
	err := r.db.QueryRow(
		"SELECT * FROM \"user\" WHERE login = $1 AND password = $2",
		login,
		hashedPassword,
	).Scan(
		&respUser.UserId,
		&respUser.FirstName,
		&respUser.LastName,
		&respUser.Login,
		&respUser.Email,
		&respUser.Password,
		&respUser.CreatedAt,
	)

	if err != nil {
		fmt.Printf("%s: %s\n", op, err)
		return models.User{}, err
	}

	return respUser, nil
}
