package handler

import (
	"net/http"

	"ps-gogo-manajer/internal/user/dto"
	"ps-gogo-manajer/internal/user/usecase"
	customErrors "ps-gogo-manajer/pkg/custom-errors"
	customValidators "ps-gogo-manajer/pkg/custom-validators"
	"ps-gogo-manajer/pkg/helper"
	"ps-gogo-manajer/pkg/jwt"
	"ps-gogo-manajer/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type UserHandler struct {
	UseCase  usecase.UserUseCase
	Validate *validator.Validate
}

func NewUserHandler(useCase usecase.UserUseCase, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		UseCase:  useCase,
		Validate: validate,
	}
}

func (c *UserHandler) AuthenticateUser(ctx echo.Context) error {
	var request = new(dto.AuthRequest)

	if err := ctx.Bind(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if err := c.Validate.Struct(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	var token *string
	var err error
	var statusCode int

	if request.Action == "create" {
		token, err = c.UseCase.Create(ctx.Request().Context(), request)
		statusCode = http.StatusCreated
	}

	if request.Action == "login" {
		token, err = c.UseCase.Login(ctx.Request().Context(), request)
		statusCode = http.StatusOK
	}

	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(statusCode, &dto.AuthResponse{
		Email:       request.Email,
		AccessToken: *token,
	})
}

func (c *UserHandler) GetUser(ctx echo.Context) error {

	userData := ctx.Get("user").(*jwt.JwtClaim)
	user, err := c.UseCase.GetUser(ctx.Request().Context(), userData.Id)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, &dto.UserResponse{
		Email:           user.Email,
		Username:        helper.DerefString(user.Username, ""),
		UserImageUri:    helper.DerefString(user.UserImageUri, ""),
		CompanyName:     helper.DerefString(user.CompanyName, ""),
		CompanyImageUri: helper.DerefString(user.CompanyImageUri, ""),
	})
}

func (c *UserHandler) UpdateUser(ctx echo.Context) error {
	var request = new(dto.UpdateUserRequest)

	if err := ctx.Bind(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if err := c.Validate.Struct(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if request.UserImageUri != nil {
		parsedImgUri, isValid := customValidators.ParseURI(*request.UserImageUri)
		if !isValid {
			err := errors.Wrap(customErrors.ErrBadRequest, "invalid image uri")
			return ctx.JSON(response.WriteErrorResponse(err))
		}
		request.UserImageUri = &parsedImgUri
	}

	if request.CompanyImageUri != nil {
		parsedImgUri, isValid := customValidators.ParseURI(*request.CompanyImageUri)
		if !isValid {
			err := errors.Wrap(customErrors.ErrBadRequest, "invalid company image uri")
			return ctx.JSON(response.WriteErrorResponse(err))
		}
		request.CompanyImageUri = &parsedImgUri
	}

	userData := ctx.Get("user").(*jwt.JwtClaim)
	user, err := c.UseCase.UpdateUser(ctx.Request().Context(), request, userData.Id)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, &dto.UserResponse{
		Email:           user.Email,
		Username:        helper.DerefString(user.Username, ""),
		UserImageUri:    helper.DerefString(user.UserImageUri, ""),
		CompanyName:     helper.DerefString(user.CompanyName, ""),
		CompanyImageUri: helper.DerefString(user.CompanyImageUri, ""),
	})
}
