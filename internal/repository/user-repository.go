package repository

import (
	"database/sql"
	"time"
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

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

type UserDatabase interface {
	CreateUser(user User) error
	FindUserById(id *int) (*User, error)
	FindUserByUsername(username *string) (*User, error)
}

func (ur UserRepository) FindUserByUsername(username *string) (*User, error) {
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

func (ur UserRepository) FindUserById(id *int) (*User, error) {
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

func (ur UserRepository) CreateUser(user *User) error {
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

func (ur UserRepository) UpdateUserDetails(user *User, id int) error {
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
