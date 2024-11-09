package request

type UserRequest struct{}

type CreateUserRequest struct {
	Name     string  `json:"name" form:"name"`
	Username string  `json:"username" form:"username"`
	Email    string  `json:"email" form:"email"`
	Password string  `json:"password" form:"password"`
	Phone    string  `json:"phone" form:"phone"`
	Address  *string `json:"address" form:"address"`
	Picture  *string `json:"picture" form:"picture"`
}
