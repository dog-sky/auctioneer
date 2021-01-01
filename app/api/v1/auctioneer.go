package v1

import (
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
)

func (h *V1Handler) SearchItemData(c *fiber.Ctx) error {
	queryParams := new(searchQueryParams)
	if err := c.QueryParser(queryParams); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Sprintf("Invalid query params. Err: %v", err),
		)
	}

	if err := checkQueryParams(queryParams); err != nil {
		return err
	}

	realmID := h.BlizzClient.GetRealmID(queryParams.RealmName)
	if realmID == 0 {
		return fiber.NewError(
			fiber.StatusNotFound,
			fmt.Sprintf("Realm %s not found", queryParams.RealmName),
		)
	}

	res := new(ResponseV1)
	res.Success = true

	// itemData
	_, err := h.BlizzClient.SearchItem(queryParams.ItemName, queryParams.Region)
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	// сделать запрос в аук по серверу
	// отдать в ответе только то, что пересекается с запрошеным предметом (или всё, если запрос пустой)

	return c.JSON(res)
}
