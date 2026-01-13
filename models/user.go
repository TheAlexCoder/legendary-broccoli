package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	RecoveryPhrase string    `json:"recovery_phrase,omitempty"`
	FiredDate      *string   `json:"fired_date"`
	CreatedAt      time.Time `json:"created_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
	IsDeleted      bool      `json:"is_deleted"`
}

func CreateUser(username, recoveryPhrase string) (*User, error) {
	result, err := DB.Exec("INSERT INTO users (username, recovery_phrase) VALUES (?, ?)",
		username, recoveryPhrase)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetUserByID(int(id))
}

func GetUserByRecoveryPhrase(recoveryPhrase string) (*User, error) {
	user := &User{}
	err := DB.QueryRow(`
		SELECT id, username, recovery_phrase, fired_date, created_at, deleted_at, is_deleted
		FROM users WHERE recovery_phrase = ? AND is_deleted = FALSE`,
		recoveryPhrase).Scan(
		&user.ID, &user.Username, &user.RecoveryPhrase, &user.FiredDate,
		&user.CreatedAt, &user.DeletedAt, &user.IsDeleted)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func GetUserByID(id int) (*User, error) {
	user := &User{}
	err := DB.QueryRow(`
		SELECT id, username, recovery_phrase, fired_date, created_at, deleted_at, is_deleted
		FROM users WHERE id = ?`,
		id).Scan(
		&user.ID, &user.Username, &user.RecoveryPhrase, &user.FiredDate,
		&user.CreatedAt, &user.DeletedAt, &user.IsDeleted)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func UpdateUser(user *User) error {
	_, err := DB.Exec(`
		UPDATE users SET username = ?, fired_date = ? WHERE id = ?`,
		user.Username, user.FiredDate, user.ID)
	return err
}

func SoftDeleteUser(id int) error {
	_, err := DB.Exec(`
		UPDATE users SET is_deleted = TRUE, deleted_at = CURRENT_TIMESTAMP WHERE id = ?`,
		id)
	return err
}

func RestoreUser(id int) error {
	_, err := DB.Exec(`
		UPDATE users SET is_deleted = FALSE, deleted_at = NULL WHERE id = ?`,
		id)
	return err
}
