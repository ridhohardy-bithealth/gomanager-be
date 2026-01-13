package config

import (
	"net/http"
	"ps-gogo-manajer/db"
	employeeHandler "ps-gogo-manajer/internal/employee/handler"
	employeeRepository "ps-gogo-manajer/internal/employee/repository"
	employeeUsecase "ps-gogo-manajer/internal/employee/usecase"

	departmentHandler "ps-gogo-manajer/internal/department/handler"
	departmentRepository "ps-gogo-manajer/internal/department/repository"
	departmentUsecase "ps-gogo-manajer/internal/department/usecase"
	fileHandler "ps-gogo-manajer/internal/files/handler"
	fileUsecase "ps-gogo-manajer/internal/files/usecase"
	auth "ps-gogo-manajer/internal/middleware"
	"ps-gogo-manajer/internal/routes"
	userHandler "ps-gogo-manajer/internal/user/handler"
	userRepository "ps-gogo-manajer/internal/user/repository"
	userUsecase "ps-gogo-manajer/internal/user/usecase"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type BootstrapConfig struct {
	App       *echo.Echo
	DB        *db.Postgres
	Log       *logrus.Logger
	Validator *validator.Validate
	S3Client  *s3.Client
}

func Bootstrap(config *BootstrapConfig) {
	employeeRepo := employeeRepository.NewEmployeeRepository(config.DB.Pool)
	employeeUseCase := employeeUsecase.NewEmployeeUsecase(*employeeRepo)
	employeeHandler := employeeHandler.NewEmployeeHandler(*employeeUseCase, config.Validator)

	userRepo := userRepository.NewUserRepository(config.DB.Pool)
	userUseCase := userUsecase.NewUserUseCase(*userRepo)
	userHandler := userHandler.NewUserHandler(*userUseCase, config.Validator)

	fileUsecase := fileUsecase.NewFileUseCase(config.S3Client)
	fileHandler := fileHandler.NewFileHandler(fileUsecase, config.Log)

	//department variable
	departmentRepo := departmentRepository.NewDepartmentRepository(config.DB.Pool)
	departmentUsecase := departmentUsecase.NewDepartmentUsecases(*departmentRepo)
	departmentHandler := departmentHandler.NewDepartmentHandler(*departmentUsecase,config.Validator)

	// * Middleware
	config.App.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	config.App.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Timeout",
		Timeout:      30 * time.Second,
	}))
	authMiddleware := auth.Auth()

	routes := routes.RouteConfig{
		App:             config.App,
		S3Client:        config.S3Client,
		EmployeeHandler: employeeHandler,
		UserHandler:     userHandler,
		AuthMiddleware:  authMiddleware,
		FileHandler:     fileHandler,
		DepartmentHandler : departmentHandler,
	}

	routes.SetupRoutes()
}
