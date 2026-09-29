package domain

import "time"

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (u User) ToUserDetailResponse() UserDetailResponse {
	return UserDetailResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}
