package repository

import (
	"database/sql"
	"fmt"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/user/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(u *domain.User) error {
	query := "INSERT INTO users (username, email, password, status) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRow(query, u.Username, u.Email, u.Password, u.Status).Scan(&u.ID)

	return err
}

func (r *UserRepository) UpdateUser(id int, u *domain.User) error {
	query := `
		UPDATE users 
		SET username = $1, 
		    email = $2, 
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING id`
	err := r.db.QueryRow(query, u.Username, u.Email, id).Scan(&u.ID)
	return err
}

func (r *UserRepository) ChangePasswordUser(id int, password string) error {
	query := `
		UPDATE users 
		SET password = $1, 
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id`
	err := r.db.QueryRow(query, password, id).Scan(&id)
	return err
}

func (r *UserRepository) ChangeStatusUser(id int, status domain.UserStatus) error {
	query := `
		UPDATE users 
		SET status = $1, 
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id
	`
	err := r.db.QueryRow(query, status, id).Scan(&id)
	return err
}

func (r *UserRepository) GetUserById(id int) (domain.User, error) {
	var u domain.User

	query := `
	SELECT 
		id, 
		username, 
		email, 
		password, 
		status, 
		created_at, 
		updated_at, 
		deleted_at 
	FROM users 
	WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, id).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)

	return u, err
}

func (r *UserRepository) GetUserByEmail(email string) (domain.User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			password,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	var u domain.User

	err := r.db.QueryRow(query, email).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("%w: %v", utils.ErrUserNotFound, err)
		}
		return u, fmt.Errorf("%w: %v", utils.ErrQueryFailed, err)
	}

	return u, nil
}

func (r *UserRepository) CheckExistsBy(column string, value any) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM users WHERE %s = $1 AND deleted_at IS NULL)", column)
	err := r.db.QueryRow(query, value).Scan(&exists)
	return exists, err
}
