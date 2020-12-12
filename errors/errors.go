package errors

import "fmt"

type AlreadyExistsError struct {
	Obj string
	Id  string
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s with ID %q already exists", e.Obj, e.Id)
}
