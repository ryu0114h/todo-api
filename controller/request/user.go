package request

type CreateUserRequestBody struct {
	Username   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Role       string `json:"role" validate:"required,oneof=admin user"`
	CompanyIds []uint `json:"company_ids" validate:"required,min=1,dive,gt=0"`
}
