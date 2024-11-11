package middlewares

import (
	"errors"
	"fiber_boilerplate/handlers"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				var errMsg string
				switch e := err.(type) {
				case error:
					errMsg = e.Error()
				default:
					errMsg = fmt.Sprintf("%v", e)
				}

				handlers.InternalServerError(c, errors.New(errMsg))
			}
		}()
		return c.Next()
	}
}
