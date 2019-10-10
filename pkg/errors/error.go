package errors

import "fmt"

type UserAlreadyExistsError struct{}

func (UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user already exists")
}

type UserNotExistsError struct{}

func (UserNotExistsError) Error() string {
	return fmt.Sprintf("user not exists")
}
