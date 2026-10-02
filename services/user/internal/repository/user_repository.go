package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"github.com/stevanusy21/golang_sandbox/services/user/internal/domain"
)

const LogLocation = "User Repository"

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(u *domain.User) error {
	query := `
		INSERT INTO users (username, email, password, status) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`

	err := r.db.QueryRow(query, u.Username, u.Email, u.Password, u.Status).Scan(&u.ID)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	return nil
}

func (r *UserRepository) UpdateUser(id int, u *domain.User) error {
	query := `
		UPDATE users 
		SET username = $1, 
		    email = $2, 
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND deleted_at IS NULL
	`

	//Tidak return rows affected karena user sudah di get sebelumnya
	if _, err := r.db.Exec(query, u.Username, u.Email, id); err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	return nil
}

func (r *UserRepository) ChangePasswordUser(id int, password string) error {
	query := `
		UPDATE users 
		SET password = $1, 
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND deleted_at IS NULL
	`
	//Tidak return rows affected karena user sudah di get sebelumnya
	if _, err := r.db.Exec(query, password, id); err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	return nil
}

func (r *UserRepository) ChangeStatusUser(id int, status domain.UserStatus) error {
	query := `
		UPDATE users 
		SET status = $1, 
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, status, id)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return utils.ErrQueryFailed
	}

	if rowsAffected == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) GetUserBy(column string, value any) (domain.User, error) {
	query := fmt.Sprintf(`
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
		WHERE %s = $1 AND deleted_at IS NULL
	`, column)

	var u domain.User

	err := r.db.QueryRow(query, value).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, utils.ErrUserNotFound
		}
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return domain.User{}, utils.ErrQueryFailed
	}

	return u, nil
}

func (r *UserRepository) CheckExistsBy(column string, value any, excludeId *int) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1 
			FROM users 
			WHERE %s = $1 AND deleted_at IS NULL`, column)

	args := []any{value}

	if excludeId != nil {
		query += " AND id != $2"
		args = append(args, *excludeId)
	}

	query += ")"

	var exists bool

	if err := r.db.QueryRow(query, args...).Scan(&exists); err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return false, fmt.Errorf("%w: %v", utils.ErrQueryFailed, err)
	}

	return exists, nil
}

func (r *UserRepository) GetAllUsers(filter domain.UserFilter) ([]domain.User, error) {
	qb := utils.NewPaginationQueryBuilder(`
		SELECT 
			id, 
			username, 
			email, 
			status,
			created_at,
			updated_at,
			deleted_at 
		FROM users 
		WHERE 1=1
	`)

	if filter.ID > 0 {
		qb.Where("AND", "id", "=", filter.ID)
	}

	if filter.Username != "" {
		qb.Where("AND", "username", "ILIKE", "%"+filter.Username+"%")
	}

	if filter.Email != "" {
		qb.Where("AND", "email", "ILIKE", "%"+filter.Email+"%")
	}

	if len(filter.Status) > 0 {
		qb.WhereAny("AND", "status", filter.Status)
	}

	qb.WhereNull("AND", "deleted_at")

	qb.OrderBy(filter.SortBy, filter.SortDir)
	qb.LimitOffset(filter.Limit, filter.Offset)

	query, args := qb.Build()

	rows, err := r.db.Query(query, args...)
	if err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return nil, utils.ErrQueryFailed
	}

	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User

		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.Status,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.DeletedAt,
		); err != nil {
			utils.LogError(LogLocation, utils.ErrScanFailed.Error(), err)
			return nil, utils.ErrScanFailed
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		utils.LogError(LogLocation, utils.ErrQueryFailed.Error(), err)
		return nil, utils.ErrQueryFailed
	}

	return users, nil
}
