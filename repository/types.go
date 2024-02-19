// This file contains types that are used in the repository layer.
package repository

type RegisterUserInput struct {
	Principal        string
	RegistrationType string
	FullName         string
	PasswordMethod   string
	PasswordSalt     string
	PasswordHash     string
}
type RegisterUserOutput struct {
	UserID           string
	Principal        string
	RegistrationType string
	FullName         string
}
type UpdateUserInput struct {
	UserID      string
	PhoneNumber *string
	FullName    *string
}
type UpdateUserOutput struct {
	PhoneNumber string
	FullName    string
}

type GetUserInput struct {
	UserID string
}
type GetUserOutput struct {
	PhoneNumber string
	FullName    string
}

type GetAuthDataInput struct {
	Principal        string
	RegistrationType string
}

type GetAuthDataOutput struct {
	UserID         string
	PasswordMethod string
	PasswordSalt   string
	PasswordHash   string
}
