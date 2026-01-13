package repository

import (
	"context"
	"ps-gogo-manajer/internal/user/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  email,
  hashed_password
) VALUES (
  $1, $2
) RETURNING id, email, hashed_password, username, user_image_uri, company_name, company_image_uri
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (r *UserRepository) CreateUser(ctx context.Context, arg CreateUserParams) (model.User, error) {
	row := r.pool.QueryRow(ctx, createUser, arg.Email, arg.HashedPassword)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.Username,
		&i.UserImageUri,
		&i.CompanyName,
		&i.CompanyImageUri,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, email, hashed_password, username, user_image_uri, company_name, company_image_uri FROM users
WHERE id = $1 LIMIT 1
`

func (r *UserRepository) GetUser(ctx context.Context, id int) (model.User, error) {
	row := r.pool.QueryRow(ctx, getUser, id)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.Username,
		&i.UserImageUri,
		&i.CompanyName,
		&i.CompanyImageUri,
	)
	return i, err
}

const getUserFromEmail = `-- name: GetUserFromEmail :one
SELECT id, email, hashed_password, username, user_image_uri, company_name, company_image_uri FROM users
WHERE email = $1 LIMIT 1
`

func (r *UserRepository) GetUserFromEmail(ctx context.Context, email string) (model.User, error) {
	row := r.pool.QueryRow(ctx, getUserFromEmail, email)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.Username,
		&i.UserImageUri,
		&i.CompanyName,
		&i.CompanyImageUri,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
  email = COALESCE($1, email),
  username = COALESCE($2, username),
  user_image_uri = COALESCE($3, user_image_uri),
  company_name = COALESCE($4, company_name),
  company_image_uri = COALESCE($5, company_image_uri)
WHERE
  id = $6
RETURNING id, email, hashed_password, username, user_image_uri, company_name, company_image_uri
`

type UpdateUserParams struct {
	Email           *string
	Username        *string
	UserImageUri    *string
	CompanyName     *string
	CompanyImageUri *string
	ID              int
}

func (r *UserRepository) UpdateUser(ctx context.Context, arg UpdateUserParams) (model.User, error) {
	row := r.pool.QueryRow(ctx, updateUser,
		arg.Email,
		arg.Username,
		arg.UserImageUri,
		arg.CompanyName,
		arg.CompanyImageUri,
		arg.ID,
	)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.Username,
		&i.UserImageUri,
		&i.CompanyName,
		&i.CompanyImageUri,
	)
	return i, err
}
