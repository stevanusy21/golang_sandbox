package domain

type UserStatus string

const (
	UserActive    UserStatus = "ACTIVE"
	UserInactive  UserStatus = "INACTIVE"
	UserSuspended UserStatus = "SUSPENDED"
)

func (s UserStatus) IsValid() bool {
	switch s {
	case UserActive, UserInactive, UserSuspended:
		return true
	default:
		return false
	}
}