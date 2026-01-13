package usecase

import (
	"context"
	"ps-gogo-manajer/internal/department/dto"
	"ps-gogo-manajer/internal/department/repository"
	customErrors "ps-gogo-manajer/pkg/custom-errors"

	"github.com/pkg/errors"
)

type DepartmentUsecase struct {
	departmentRepo repository.DepartmentRepository
}

func NewDepartmentUsecases(departmentRepo repository.DepartmentRepository) *DepartmentUsecase {
	return &DepartmentUsecase{
		departmentRepo: departmentRepo,
	}
}

func (u *DepartmentUsecase) CreateDepartment(ctx context.Context, userID int, payload *dto.CreateDepartmentPayload) (*dto.Department, error) {
	return u.departmentRepo.CreateDepartment(ctx, userID, payload)
}

func (u *DepartmentUsecase) GetListDepartment(ctx context.Context, userID int, payload *dto.GetDepartmentListParams) (*[]dto.Department, error) {
	return u.departmentRepo.GetListDepartment(ctx, userID, payload)
}

func (u *DepartmentUsecase) UpdateDepartment(ctx context.Context, userID int, departmentId int, payload *dto.PatchDepartmentPayload) (*dto.Department, error) {
	isDepartmentExists, err := u.departmentRepo.CheckIfDepartmentExist(ctx, userID, departmentId)
	if err != nil {
		return nil, err
	}
	if !isDepartmentExists {
		return nil, errors.Wrap(customErrors.ErrNotFound, "department id for this user not found")
	}
	return u.departmentRepo.UpdateDepartment(ctx, userID, departmentId, payload)
}

func (u *DepartmentUsecase) DeleteDepartment(ctx context.Context, userID int, departmentID int) error {

	isDepartmentExists, err := u.departmentRepo.CheckIfDepartmentExist(ctx, userID, departmentID)
	if err != nil {
		return err
	}

	if !isDepartmentExists {
		return errors.Wrap(customErrors.ErrNotFound, "department not exist")
	}

	isEmployeeExists, err := u.departmentRepo.CheckIfEmployeeExist(ctx, userID, departmentID)
	if err != nil {
		return err
	}

	if isEmployeeExists {
		return errors.Wrap(customErrors.ErrConflict, "still containing employee")
	}

	return u.departmentRepo.DeleteDepartment(ctx, userID, departmentID)
}
