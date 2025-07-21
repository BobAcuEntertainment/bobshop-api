package domain

type Role int

const (
	UserRole Role = iota
	AdminRole
)

var roleNames = []string{"user", "admin"}

func (r Role) String() string {
	return roleNames[r]
}
