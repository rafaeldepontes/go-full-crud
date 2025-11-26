package repository

import (
	"database/sql"
	"time"

	"github.com/rafaeldepontes/go-full-crud/internal/util"
)

type User struct {
	Id        *int       `json:"id"`
	Username  *string    `json:"username"`
	Password  *string    `json:"password"`
	Email     *string    `json:"email,omitempty"`
	Birthdate *string    `json:"birthdate,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (ur *Repository) SetDb(db *sql.DB) { ur.db = db }

func (ur *Repository) Ping() bool {
	db := ur.db
	if db == nil {
		return false
	}
	return db.Ping() == nil
}

type UserDatabase interface {
	CreateUser(user User) error
	FindUserById(id *int) (*User, error)
	FindUserByUsername(username *string) (*User, error)
}

func (ur *Repository) FindUserByUsername(username *string) (*User, error) {
	user := User{}

	query :=
		`
		SELECT id, username, password, email, birthdate, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	err := ur.db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Birthdate, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *Repository) FindUserById(id *int) (*User, error) {
	user := User{}

	query :=
		`
		SELECT id, username, password, email, birthdate, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := ur.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Birthdate, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *Repository) CreateUser(user *User) error {
	query :=
		`
	INSERT INTO users (username, password, created_at)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	err := ur.db.QueryRow(query, user.Username, user.Password, user.CreatedAt).Scan(&user.Id, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (ur *Repository) UpdateUserDetails(user *User, id int) error {
	query :=
		`
	UPDATE users 
	SET email = $1, birthdate = $2, updated_at = $3
	WHERE id = $4
	RETURNING email, birthdate, updated_at
	`

	err := ur.db.QueryRow(query, user.Email, user.Birthdate, user.UpdatedAt, id).Scan(&user.Email, &user.Birthdate, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (ur *Repository) DeleteUserById(id int) error {
	if id == 0 {
		return util.ErrorBlankId
	}

	query :=
		`
	DELETE FROM users
	where id = $1
	`

	ur.db.QueryRow(query, id)

	return nil
}
