package uuid

import (
	uuid "github.com/iris-contrib/go.uuid"
)

func Get() string {
	u, _ := uuid.NewV4()
	return u.String()
}
