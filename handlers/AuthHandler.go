package handlers

import (
	"fiber_boilerplate/helpers"
	"fiber_boilerplate/models"
	"fiber_boilerplate/types/requests"
	"fiber_boilerplate/types/responses"
	validators "fiber_boilerplate/validator"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct{}

var authModel = new(models.AuthModel)
var authUserModel = new(models.UserModel)
var otpModel = new(models.OtpModel)

func (a *AuthHandler) Refresh(ctx *fiber.Ctx) error {
	var r requests.RefreshTokenRequest

	if err := ctx.BodyParser(&r); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse(ctx, err))
	}

	if err := validators.ValidateStruct(&r); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse(ctx, err))
	}

	token, err := jwt.Parse(r.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return ErrorResponse(ctx, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ErrorResponse(ctx, err)
	}

	if token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string)
		if !ok {
			return ErrorResponse(ctx, err)
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return ErrorResponse(ctx, err)
		}
		deleted, delErr := authModel.DestroyAuth(refreshUUID)
		if delErr != nil || deleted == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid authorization, please login again",
			})
		}

		ts, createErr := authModel.GenerateToken(userID)
		if createErr != nil {
			return ErrorResponse(ctx, createErr)
		}

		saveErr := authModel.StartAuth(userID, ts)
		if saveErr != nil {
			return ErrorResponse(ctx, saveErr)
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return ctx.Status(fiber.StatusOK).JSON(responses.APIResponse{
			Success: true,
			Message: "Token refreshed successfully",
			Data:    fiber.Map{"tokens": tokens},
			Status:  fiber.StatusOK,
		})
	} else {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid authorization, please login again",
		})
	}
}

func (a *AuthHandler) Login(ctx *fiber.Ctx) error {
	var r requests.LoginUserRequest

	if err := ctx.BodyParser(&r); err != nil {
		return ErrorResponse(ctx, err)
	}

	if err := validators.ValidateStruct(&r); err != nil {
		return ErrorResponse(ctx, err)
	}

	user, err := authUserModel.Login(r.Email)
	if err != nil {
		return ErrorResponse(ctx, err)
	}

	err = helpers.ComparePassword(user.Password, r.Password)
	if err != nil {
		return ErrorResponse(ctx, err)
	}

	var token models.Token

	tokenDetails, err := authModel.GenerateToken(user.ID)
	if err != nil {
		return ErrorResponse(ctx, err)
	}

	saveErr := authModel.StartAuth(user.ID, tokenDetails)
	if saveErr == nil {
		token.AccessToken = tokenDetails.AccessToken
		token.RefreshToken = tokenDetails.RefreshToken
	}

	if user.TwoFactorAuth {
		otp := models.OTP{
			UserID:     user.ID,
			OtpCode:    otpModel.GenerateRandomNumber(),
			IsVerified: 0,
			IsExpired:  0,
		}

		_, err = otpModel.Create(&otp)

		if err != nil {
			return ErrorResponse(ctx, err)
		}

		ctx.Status(fiber.StatusOK).JSON(responses.APIResponse{
			Success: true,
			Message: "User logged in successfully",
			Data: fiber.Map{
				"user":  user,
				"2fa":   true,
				"phone": user.Phone,
				"otp":   otp.OtpCode,
			},
			Status: fiber.StatusOK,
		})

	} else {
		ctx.Status(fiber.StatusOK).JSON(responses.APIResponse{
			Success: true,
			Message: "User logged in successfully",
			Data: fiber.Map{
				"2fa":   false,
				"user":  user,
				"token": token,
			},
			Status: fiber.StatusOK,
		})
	}

	return nil
}

func (a *AuthHandler) Register(ctx *fiber.Ctx) error {
	var r requests.CreateUserRequest

	if err := ctx.BodyParser(&r); err != nil {
		return ErrorResponse(ctx, err)

	}

	if err := validators.ValidateStruct(&r); err != nil {
		return ErrorResponse(ctx, err)
	}

	hashedPassword, err := helpers.HashPassword(string(r.Password))
	if err != nil {
		return ErrorResponse(ctx, err)
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
		return ErrorResponse(ctx, err)
	}

	_, err = authUserModel.FindByUsername(r.Username)
	if err == nil {
		return ErrorResponse(ctx, err)
	}

	_, err = authUserModel.FindByPhone(*r.Phone)
	if err == nil {
		return ErrorResponse(ctx, err)
	}

	user, err := authUserModel.Create(&newUser)
	if err != nil {
		return ErrorResponse(ctx, err)
	}

	otp := models.OTP{
		UserID:     user.ID,
		OtpCode:    otpModel.GenerateRandomNumber(),
		IsVerified: 0,
		IsExpired:  0,
	}

	_, err = otpModel.Create(&otp)
	if err != nil {
		return ErrorResponse(ctx, err)
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
		return ErrorResponse(ctx, err)
	}

	user, err := authUserModel.FindByPhone(r.PhoneNumber)
	if err != nil {
		return ErrorResponse(ctx, err)
	}

	findOtp, err := otpModel.FindOtp(user.ID, r.Otp)
	if err != nil {
		return ErrorResponse(ctx, err)
	}

	fmt.Println(findOtp)

	user.IsActive = true
	_, err = authUserModel.Update(user)
	if err != nil {
		return ErrorResponse(ctx, err)
	}
	otpModel.Delete(findOtp.ID)

	var token models.Token

	tokenDetails, err := authModel.GenerateToken(user.ID)
	if err != nil {
		return ErrorResponse(ctx, err)
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
