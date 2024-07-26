package common

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrTableCreationFailed = errors.New("unable to create table")
	ErrInternalServerError = errors.New("internal server error")
	ErrRecordNotFound      = errors.New("not found")
	ErrDuplicateKey        = errors.New("already exists")
	ErrInvalidID           = errors.New("invalid ID")
	ErrBadRequest          = errors.New("bad request")
	ErrNotImplemented      = errors.New("not implemented")
	ErrReadingRequestBody  = errors.New("Error reading response body")
	ErrExpiredSubscription = errors.New("Subscription expired. Please renew subscription.")
	ErrQueryLimit          = errors.New("query limit exceeded")
	ErrReadingResponse     = errors.New("err reading response body")
	ErrNoSubscription      = errors.New("no subscription exists")
	ErrUnauthorized        = errors.New("Unauthorized")
)

func HandleError(c *fiber.Ctx, err error) error {
	switch err {
	case ErrBadRequest:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	case fiber.ErrForbidden: // TODO: refactor
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	case ErrRecordNotFound:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": ErrRecordNotFound.Error()})
	case ErrNoSubscription:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": ErrNoSubscription.Error()})
	case ErrDuplicateKey:
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": ErrDuplicateKey.Error()})
	case ErrUnauthorized:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": ErrUnauthorized.Error()})
	case ErrInvalidID:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrInvalidID.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": ErrInternalServerError.Error()})
	}
}
