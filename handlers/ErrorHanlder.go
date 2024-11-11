package handlers

import (
	"fiber_boilerplate/types/responses"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ErrorResponse(c *fiber.Ctx, err error) error {
	if fiberErr, ok := err.(*fiber.Error); ok {
		return c.Status(fiberErr.Code).JSON(responses.APIResponse{
			Message: fiberErr.Message,
			Error:   err.Error(),
		})
	}

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, fieldErr := range validationErrs {
			errorMessages = append(errorMessages, fieldErr.Error())
		}

		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{
			Message: "Bad Request: Invalid input data",
			Error:   strings.Join(errorMessages, ", "),
		})
	}

	log.Printf("Internal Server Error: %v", err)
	return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{
		Message: "Terjadi kesalahan di server.",
		Error:   err.Error(),
	})
}

func InternalServerError(c *fiber.Ctx, err error) error {
	log.Printf("Internal Server Error: %v", err)
	return c.Status(fiber.StatusInternalServerError).JSON(responses.APIResponse{
		Success: false,
		Status:  500,
		Message: "Something went wrong.",
		Error:   err.Error(),
	})
}
