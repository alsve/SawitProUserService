package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	stderr "github.com/alsve/SawitProUserService/internal/error"
)

// RegisterUser register user.
func (r *Repository) RegisterUser(ctx context.Context, input RegisterUserInput) (output *RegisterUserOutput, err error) {
	subCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	tx, err := r.Db.BeginTx(subCtx, &sql.TxOptions{Isolation: sql.LevelDefault})
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(subCtx,
		`INSERT INTO public."user"(id, phone_number, registration_type, full_name) VALUES ($1, $2, $3, $4) RETURNING id, phone_number, registration_type, full_name`,
	)
	if err != nil {
		return nil, stderr.NewServerError(stderr.ErrInsertUserSQLStmt, err)
	}
	defer stmt.Close()

	userId := uuid.New()
	row := stmt.QueryRowContext(subCtx, userId, input.Principal, input.RegistrationType, input.FullName)

	output = &RegisterUserOutput{}
	if err := row.Scan(&output.UserID, &output.Principal, &output.RegistrationType, &output.FullName); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, stderr.NewConflictRequest("phone number has been registered")
		}
		return nil, stderr.NewServerError(stderr.ErrInsertUserScan, err)
	}

	stmt, err = tx.PrepareContext(subCtx,
		`INSERT INTO public."user_credential" (user_id, method, hash, salt) VALUES ($1, $2, $3, $4)`,
	)
	if err != nil {
		return nil, stderr.NewServerError(stderr.ErrInsertUserSQLStmt, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(subCtx, userId, input.PasswordMethod, input.PasswordHash, input.PasswordSalt)
	if err != nil {
		return nil, stderr.NewServerError(stderr.ErrInsertUserSQLStmt, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, stderr.NewServerError(stderr.ErrInsertUserCommit, err)
	}

	return output, nil
}

// UpdateUser update user.
func (r *Repository) UpdateUser(ctx context.Context, input UpdateUserInput) (output *UpdateUserOutput, err error) {
	subCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// build set fragments.
	setFragments := ""
	bindArgs := []interface{}{}
	count := 0
	if input.FullName != nil {
		setFragments += fmt.Sprintf("full_name = $%d", count+1)
		bindArgs = append(bindArgs, *input.FullName)
		count++
	}
	if input.PhoneNumber != nil {
		if count > 0 {
			setFragments += ", "
		}
		setFragments += fmt.Sprintf("phone_number = $%d", count+1)
		bindArgs = append(bindArgs, *input.PhoneNumber)
		count++
	}
	if count < 1 {
		return nil, stderr.NewBadRequest("please provide what to update")
	}

	bindArgs = append(bindArgs, input.UserID)
	stmt, err := r.Db.PrepareContext(subCtx,
		fmt.Sprintf(`UPDATE public."user" SET %s WHERE id = $%d RETURNING phone_number, full_name`, setFragments, count+1),
	)
	if err != nil {
		return nil, stderr.NewServerError(stderr.ErrGetPasswordByUserIDSQLStmt, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(subCtx, bindArgs...)

	output = &UpdateUserOutput{}
	if err := row.Scan(&output.PhoneNumber, &output.FullName); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, stderr.NewConflictRequest("phone number has been registered")
		}
		return nil, stderr.NewServerError(stderr.ErrGetPasswordByUserIDScan, err)
	}

	return output, nil
}

// GetUser retrieves user.
func (r *Repository) GetUser(ctx context.Context, input GetUserInput) (output *GetUserOutput, err error) {
	subCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	stmt, err := r.Db.PrepareContext(subCtx,
		`SELECT phone_number, full_name FROM public."user" WHERE id = $1`,
	)
	if err != nil {
		return nil, stderr.NewServerError(stderr.ErrGetPasswordByUserIDSQLStmt, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(subCtx, &input.UserID)

	output = &GetUserOutput{}
	if err := row.Scan(&output.PhoneNumber, &output.FullName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, stderr.NewBadRequest("no such user")
		}
		return nil, stderr.NewServerError(stderr.ErrGetPasswordByUserIDScan, err)
	}

	return output, nil
}

// GetAuthData retrieves password and user id.
func (r *Repository) GetAuthData(ctx context.Context, input GetAuthDataInput) (output *GetAuthDataOutput, err error) {
	subCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	stmt, err := r.Db.PrepareContext(subCtx,
		`SElECT u.id, uc.method, uc.hash, uc.salt FROM public."user" u JOIN public."user_credential" uc ON u.id = uc.user_id WHERE uc.method = 'bcrypt' and u.phone_number = $1`,
	)
	if err != nil {
		return nil, stderr.NewServerError(stderr.ErrGetPasswordByUserIDSQLStmt, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(subCtx, &input.Principal)

	output = &GetAuthDataOutput{}
	if err := row.Scan(&output.UserID, &output.PasswordMethod, &output.PasswordSalt, &output.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, stderr.NewBadRequest("wrong user/password")
		}
		return nil, stderr.NewServerError(stderr.ErrGetPasswordByUserIDScan, err)
	}

	return output, nil
}
