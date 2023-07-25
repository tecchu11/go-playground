package auth

// Role for performing API.
type Role int

const (
	ADMIN Role = iota
	USER
)

var roleMap = map[Role]string{
	ADMIN: "ADMIN",
	USER:  "USER",
}

func (role Role) String() string {
	return roleMap[role]
}
