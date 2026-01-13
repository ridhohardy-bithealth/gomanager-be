package repository

import (
	"context"
	"ps-gogo-manajer/internal/department/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type DepartmentRepository struct {
	pool *pgxpool.Pool
}

func NewDepartmentRepository(pool *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{pool: pool}
}

const (
	queryCreateDepartment = `
	INSERT INTO departments
	(
		name,
		user_id
	)
	VALUES (@name,@userID)
	RETURNING id,name`

	queryGetListDepartment = `
	SELECT
		id,
		name
	FROM departments
	WHERE
		user_id = @userID
		AND (NULLIF(@name, '') is NULL OR name ILIKE '%' || NULLIF(@name, '') || '%' )
	OFFSET @offset
	LIMIT @limit;`

	queryUpdateDepartment = `
	WITH
	payload as (
		SELECT
			NULLIF(t.name,'') AS name
		FROM(
			VALUES (
				@name
				)
		)AS t(
				name
			)
	)
	UPDATE departments
	SET
		name = COALESCE(payload.name, departments.name)
	FROM payload
	WHERE
		departments.id = @id
	RETURNING
	departments.id,
	departments.name;`

	queryCheckIsDepartmentExist = `
	SELECT EXISTS (
		SELECT id
		FROM departments
		WHERE
			user_id = @userID
			AND id = NULLIF(@id, 0)::bigint
	) is_exists;`

	queryDeleteDepartment = `
	DELETE FROM departments WHERE id = @departmentId;
	`
	queryCheckIfEmployeeExists = `
	SELECT EXISTS (
		SELECT id
		FROM employees
		WHERE
			user_id = @userID
			AND id = NULLIF(@id, 0)::bigint
	) is_exists;`
)

func (r *DepartmentRepository) CreateDepartment(ctx context.Context, userID int, payload *dto.CreateDepartmentPayload) (*dto.Department, error) {
	var department dto.Department

	args := pgx.NamedArgs{
		"name":   payload.Name,
		"userID": userID,
	}

	err := r.pool.QueryRow(ctx, queryCreateDepartment, args).Scan(
		&department.DepartmentId,
		&department.Name,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create department")
	}

	return &department, nil
}

func (r *DepartmentRepository) GetListDepartment(ctx context.Context, userID int, payload *dto.GetDepartmentListParams) (*[]dto.Department, error) {

	var departments []dto.Department

	args := pgx.NamedArgs{
		"userID": userID,
		"name":   payload.Name,
		"limit":  payload.Limit,
		"offset": payload.Offset,
	}

	rows, err := r.pool.Query(ctx, queryGetListDepartment, args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list department")
	}

	for rows.Next() {
		department := dto.Department{}
		err := rows.Scan(
			&department.DepartmentId,
			&department.Name,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse sql response")
		}
		departments = append(departments, department)
	}

	return &departments, nil
}

func (r *DepartmentRepository) UpdateDepartment(ctx context.Context, userID int, departmentId int, payload *dto.PatchDepartmentPayload) (*dto.Department, error) {

	var department dto.Department

	args := pgx.NamedArgs{
		"name": payload.Name,
		"id":   departmentId,
	}

	err := r.pool.QueryRow(ctx, queryUpdateDepartment, args).Scan(
		&department.DepartmentId,
		&department.Name,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to update departments")
	}

	return &department, nil
}

func (r *DepartmentRepository) CheckIfDepartmentExist(ctx context.Context, UserID int, departmentID int) (bool, error) {
	var isExist bool
	args := pgx.NamedArgs{
		"userID": UserID,
		"id":     departmentID,
	}
	err := r.pool.QueryRow(ctx, queryCheckIsDepartmentExist, args).Scan(&isExist)
	if err != nil {
		return false, errors.Wrap(err, "failed to check is department exists")
	}
	return isExist, nil
}

func (r *DepartmentRepository) CheckIfEmployeeExist(ctx context.Context, UserID int, departmentID int) (bool, error) {
	var isExist bool
	args := pgx.NamedArgs{
		"userID": UserID,
	}
	err := r.pool.QueryRow(ctx, queryCheckIsDepartmentExist, args).Scan(&isExist)
	if err != nil {
		return false, errors.Wrap(err, "failed to check is employee exists")
	}
	return isExist, nil
}

func (r *DepartmentRepository) DeleteDepartment(ctx context.Context, userID int, departmentID int) error {
	args := pgx.NamedArgs{
		"userID":       userID,
		"departmentId": departmentID,
	}

	_, err := r.pool.Exec(ctx, queryDeleteDepartment, args)
	if err != nil {
		return errors.Wrap(err, "failed to delete department")
	}
	return nil
}
