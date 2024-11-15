package requests

type UserRequest struct{}

type CreateUserRequest struct {
	Name     string  `json:"name" form:"name" binding:"required"`
	Username string  `json:"username" form:"username" binding:"required"`
	Email    string  `json:"email" form:"email" binding:"required" validate:"email"`
	Password string  `json:"password" form:"password" validate:"required,min=5"`
	Phone    *string `json:"phone" form:"phone"`
	Address  *string `json:"address" form:"address"`
	Picture  *string `json:"picture" form:"picture"`
}

type LoginUserRequest struct {
	Email    string `json:"email" form:"email" binding:"required" validate:"email"`
	Password string `json:"password" form:"password" binding:"required" validate:"min=5"`
}
