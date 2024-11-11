package requests

type OtpRequest struct{}

type VerifyOtpRequest struct {
	PhoneNumber string `json:"phone_number" form:"phone_number" binding:"required"`
	Otp         string `json:"otp" form:"otp" binding:"required"`
}
