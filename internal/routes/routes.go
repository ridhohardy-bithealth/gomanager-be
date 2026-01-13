package routes

import (
	"net/http"
	departmentHandler "ps-gogo-manajer/internal/department/handler"
	employeeHandler "ps-gogo-manajer/internal/employee/handler"
	fileHandler "ps-gogo-manajer/internal/files/handler"
	userHandler "ps-gogo-manajer/internal/user/handler"
	"ps-gogo-manajer/pkg/response"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App               *echo.Echo
	S3Client          *s3.Client
	FileHandler       *fileHandler.FileHandler
	EmployeeHandler   *employeeHandler.EmployeeHandler
	UserHandler       *userHandler.UserHandler
	AuthMiddleware    echo.MiddlewareFunc
	DepartmentHandler *departmentHandler.DepartmentHandler
}

func (r *RouteConfig) SetupRoutes() {
	r.setupPublicRoutes()
	r.setupAuthRoutes()
}

func (r *RouteConfig) setupPublicRoutes() {
	r.App.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.BaseResponse{
			Status:  "Ok",
			Message: "",
		})
	})

	r.App.POST("/v1/auth", r.UserHandler.AuthenticateUser)
}
func (r *RouteConfig) setupAuthRoutes() {
	v1 := r.App.Group("/v1")

	r.setupEmployeeRoute(v1)
	r.setupUserRoute(v1)
	r.setupFileRoutes(v1)
	r.setupDepartmentRoute(v1)
}

func (r *RouteConfig) setupEmployeeRoute(api *echo.Group) {
	employee := api.Group("/employee", r.AuthMiddleware)
	employee.GET("", r.EmployeeHandler.GetListEmployee)
	employee.POST("", r.EmployeeHandler.CreateEmployee)
	employee.PATCH("/:identityNumber", r.EmployeeHandler.UpdateEmployee)
	employee.DELETE("/:identityNumber", r.EmployeeHandler.DeleteEmployee)
}

func (r *RouteConfig) setupUserRoute(api *echo.Group) {
	user := api.Group("/user", r.AuthMiddleware)
	user.GET("", r.UserHandler.GetUser)
	user.PATCH("", r.UserHandler.UpdateUser)
}

func (r *RouteConfig) setupFileRoutes(api *echo.Group) {
	api.POST("/file", r.FileHandler.UploadFile, r.AuthMiddleware)
}

func (r *RouteConfig) setupDepartmentRoute(api *echo.Group) {
	department := api.Group("/department", r.AuthMiddleware)

	department.GET("", r.DepartmentHandler.GetListDepartment)
	department.POST("", r.DepartmentHandler.CreateDepartment)
	department.PATCH("/:departmentId", r.DepartmentHandler.UpdateDepartment)
	department.DELETE("/:departmentId", r.DepartmentHandler.DeleteDepartment)
}
