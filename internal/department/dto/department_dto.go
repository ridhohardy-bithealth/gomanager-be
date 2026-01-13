package dto

type Department struct {
	DepartmentId string `json:"departmentId"`
	Name         string `json:"name"`
}

type CreateDepartmentPayload struct {
	Name string `json:"name" validate:"required,min=4,max=33"`
}

type GetDepartmentListParams struct {
	Limit  int
	Offset int
	Name   string `query:"name"`
}

type PatchDepartmentPayload struct {
	Name string `json:"name" validate:"required,omitempty,min=4,max=33"`
}

type UpdateDeletePathParam struct {
	IdentityNumber string `param:"identityNumber" validate:"required"`
}
