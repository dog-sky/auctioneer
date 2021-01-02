package v1

import (
	"github.com/gofiber/fiber/v2"
)

func checkQueryParams(q *searchQueryParams) error {
	e := &fiber.Error{
		Code: fiber.StatusBadRequest,
	}

	if q.RealmName == "" {
		e.Message = "Realm name must not be empty"
		return e
	}
	if q.Region == "" {
		e.Message = "Region must not be empty"
		return e
	}
	if q.ItemName == "" {
		e.Message = "Item name must not be empty"
		return e
	}

	return nil
}
