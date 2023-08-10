package preauth

// Role for performing API.
type Role int

const (
	UNDEFINED Role = iota
	ADMIN
	USER
)

var mapByRole = map[Role]string{
	UNDEFINED: "",
	ADMIN:     "ADMIN",
	USER:      "USER",
}

var mapByLiteral = map[string]Role{
	"ADMIN": ADMIN,
	"USER":  USER,
}

func (role Role) String() string {
	return mapByRole[role]
}

func RoleFrom(literal string) Role {
	v, ok := mapByLiteral[literal]
	if !ok {
		return UNDEFINED
	}
	return v
}
