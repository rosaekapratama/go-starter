package keycloak

import (
	"fmt"
	"github.com/rosaekapratama/go-starter/constant/str"
	"strings"
)

func SplitFullName(fullName string) (firstName string, lastName string) {
	names := strings.SplitN(strings.TrimSpace(fullName), str.Space, 2)
	if len(names) > 1 {
		return names[0], names[1]
	} else {
		return names[0], str.Empty
	}
}

func ConcatFirstNameLastName(firstName string, lastName string) (fullName string) {
	return strings.TrimSpace(fmt.Sprintf("%s %s", strings.TrimSpace(firstName), strings.TrimSpace(lastName)))
}
