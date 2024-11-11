package handlers

import (
	"fiber_boilerplate/helpers"
	"fiber_boilerplate/models"
	"fiber_boilerplate/types/requests"
	"fiber_boilerplate/types/responses"
	validators "fiber_boilerplate/validator"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct{}

var authModel = new(models.AuthModel)
var authUserModel = new(models.UserModel)
var otpModel = new(models.OtpModel)

func (a *AuthHandler) Refresh(ctx *fiber.Ctx) {

}

func (a *AuthHandler) Login(ctx *fiber.Ctx) {

}

func (a *AuthHandler) Register(ctx *fiber.Ctx) error {
	var r requests.CreateUserRequest

	if err := ctx.BodyParser(&r); err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})

	}

	if err := validators.ValidateStruct(&r); err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	hashedPassword, err := helpers.HashPassword(string(r.Password))
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
		return err
	}

	newUser := models.User{
		Email:         r.Email,
		Password:      hashedPassword,
		Role:          "user",
		Name:          r.Name,
		Username:      r.Username,
		Phone:         r.Phone,
		Address:       r.Address,
		Picture:       r.Picture,
		IsActive:      false,
		TwoFactorAuth: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = authUserModel.FindByEmail(r.Email)
	if err == nil {
		ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
		return err
	}

	_, err = authUserModel.FindByUsername(r.Username)
	if err == nil {
		ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Username already exists",
		})
		return err
	}

	_, err = authUserModel.FindByPhone(*r.Phone)
	if err == nil {
		ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Phone number already exists",
		})
		return err
	}

	user, err := authUserModel.Create(&newUser)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
		return err
	}

	otp := models.OTP{
		UserID:     user.ID,
		OtpCode:    otpModel.GenerateRandomNumber(),
		IsVerified: 0,
		IsExpired:  0,
	}

	_, err = otpModel.Create(&otp)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create otp",
		})
		return err
	}

	ctx.Status(fiber.StatusCreated).JSON(responses.APIResponse{
		Success: true,
		Message: "User created successfully",
		Data: fiber.Map{
			"user": user,
			"otp":  otp.OtpCode,
		},
		Status: fiber.StatusCreated,
	})

	return nil
}

func (a *AuthHandler) VerifyOtp(ctx *fiber.Ctx) error {
	var r requests.VerifyOtpRequest

	if err := ctx.BodyParser(&r); err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	user, err := authUserModel.FindByPhone(r.PhoneNumber)
	if err != nil {
		ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	findOtp, err := otpModel.FindOtp(user.ID, r.Otp)
	if err != nil {
		ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Otp not found",
		})
	}

	fmt.Println(findOtp)

	user.IsActive = true
	_, err = authUserModel.Update(user)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to activate user",
		})
	}
	otpModel.Delete(findOtp.ID)

	var token models.Token

	tokenDetails, err := authModel.GenerateToken(user.ID)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
		})
	}

	saveErr := authModel.StartAuth(user.ID, tokenDetails)
	if saveErr == nil {
		token.AccessToken = tokenDetails.AccessToken
		token.RefreshToken = tokenDetails.RefreshToken
	}

	ctx.Status(fiber.StatusOK).JSON(responses.APIResponse{
		Success: true,
		Message: "User activated successfully",
		Data: fiber.Map{
			"user":  user,
			"token": token,
		},
		Status: fiber.StatusOK,
	})

	return nil
}
