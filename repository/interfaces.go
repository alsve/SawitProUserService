// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	RegisterUser(ctx context.Context, input RegisterUserInput) (output *RegisterUserOutput, err error)
	UpdateUser(ctx context.Context, input UpdateUserInput) (output *UpdateUserOutput, err error)
	GetUser(ctx context.Context, input GetUserInput) (output *GetUserOutput, err error)
	GetAuthData(ctx context.Context, input GetAuthDataInput) (output *GetAuthDataOutput, err error)
}
