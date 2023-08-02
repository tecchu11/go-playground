package preauth

import (
	"errors"
	"fmt"
)

type AuthorizationManager interface {
	Authorize(Role) error
}

// AuthiorizedList roles to perform API.
type AuthorizedList []Role

func NewAuthorizationManager(permittedRoles []Role) AuthorizationManager {
	var list AuthorizedList = permittedRoles
	return list
}

// Authorize with passsed role by AuthorizedList
func (list AuthorizedList) Authorize(role Role) error {
	if role == UNDIFINED {
		return errors.New("authorization failuer because role is undifined")
	}
	for _, v := range list {
		if v == role {
			return nil
		}
	}
	return fmt.Errorf("authorization failuer becasue permited role are %v but requested user having role %s", list, role.String())
}
