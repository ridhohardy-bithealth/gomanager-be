package handler

import (
	"net/http"
	"ps-gogo-manajer/internal/department/dto"
	"ps-gogo-manajer/internal/department/usecase"
	customErrors "ps-gogo-manajer/pkg/custom-errors"
	customValidators "ps-gogo-manajer/pkg/custom-validators"
	"ps-gogo-manajer/pkg/jwt"
	"ps-gogo-manajer/pkg/response"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type DepartmentHandler struct {
	departmentUsecase usecase.DepartmentUsecase
	validator         *validator.Validate
}

const (
	DEFAULT_LIMIT  = 5
	DEFAULT_OFFSET = 0
)

func NewDepartmentHandler(department usecase.DepartmentUsecase, validator *validator.Validate) *DepartmentHandler {
	return &DepartmentHandler{
		departmentUsecase: department,
		validator:         validator,
	}
}

func (h DepartmentHandler) CreateDepartment(ctx echo.Context) error {
	var payload dto.CreateDepartmentPayload

	if err := ctx.Bind(&payload); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if err := h.validator.Struct(payload); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	userData := ctx.Get("user").(*jwt.JwtClaim)

	department, err := h.departmentUsecase.CreateDepartment(ctx.Request().Context(), userData.Id, &payload)

	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusCreated, department)
}

func (h DepartmentHandler) GetListDepartment(ctx echo.Context) error {

	limitStr := ctx.QueryParam("limit")
	offsetStr := ctx.QueryParam("offset")

	limit := customValidators.ParseLimitOffset(limitStr, DEFAULT_LIMIT)
	offset := customValidators.ParseLimitOffset(offsetStr, DEFAULT_OFFSET)

	payload := dto.GetDepartmentListParams{
		Limit:  limit,
		Offset: offset,
	}

	if err := ctx.Bind(&payload); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	userData := ctx.Get("user").(*jwt.JwtClaim)

	departments, err := h.departmentUsecase.GetListDepartment(ctx.Request().Context(), userData.Id, &payload)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if len(*departments) == 0 {
		return ctx.JSON(http.StatusOK, make([]string, 0))
	}

	return ctx.JSON(http.StatusOK, &departments)
}

func (h DepartmentHandler) UpdateDepartment(ctx echo.Context) error {
	departmentId := ctx.Param("departmentId")

	if departmentId == "" {
		err := errors.Wrap(customErrors.ErrBadRequest, "department id required")
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	var payload dto.PatchDepartmentPayload

	if err := ctx.Bind(&payload); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if err := h.validator.Struct(payload); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	id, err := strconv.Atoi(departmentId)
	if err != nil {
		err = errors.Wrap(customErrors.ErrNotFound, "department id not found")
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	userData := ctx.Get("user").(*jwt.JwtClaim)

	department, err := h.departmentUsecase.UpdateDepartment(ctx.Request().Context(), userData.Id, id, &payload)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, department)
}

func (h DepartmentHandler) DeleteDepartment(ctx echo.Context) error {

	departmentId := ctx.Param("departmentId")

	if departmentId == "" {
		err := errors.Wrap(customErrors.ErrBadRequest, "department id required")
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	id, err := strconv.Atoi(departmentId)
	if err != nil {
		err = errors.Wrap(customErrors.ErrNotFound, "wrong department id")
		return ctx.JSON(response.WriteErrorResponse(err))
	}
	userData := ctx.Get("user").(*jwt.JwtClaim)

	err = h.departmentUsecase.DeleteDepartment(ctx.Request().Context(), userData.Id, id)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  http.StatusText(http.StatusOK),
		Message: "deleted",
	})

}
