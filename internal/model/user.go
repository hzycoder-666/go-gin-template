package model

import "time"

type Role string

type User struct {
	ID        int64
	Username  string
	Password  string
	Nickname  *string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	RoleGuest  Role = "guest"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

func IsAdmin(r Role) bool {
	return r == RoleAdmin
}

func IsValid(r Role) bool {
	switch r {
	case RoleGuest, RoleMember, RoleAdmin:
		return true
	default:
		return false
	}
}
