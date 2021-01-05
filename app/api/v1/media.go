package v1

import (
	// "auctioneer/app/blizz"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func (h *V1Handler) SearchItemMedia(c *fiber.Ctx) error {
	itemID := c.Params("item_id")

	res, err := h.blizzClient.GetItemMedia(itemID)
	if err != nil {
		if err != nil {
			return fiber.NewError(
				fiber.StatusBadRequest,
				err.Error(),
			)
		}
	}

	if res == nil {
		return fiber.NewError(
			fiber.StatusNotFound,
			fmt.Sprintf("Item with ID %s not found", itemID),
		)
	}

	resp := new(ResponseV1ItemMedia)
	resp.Success = true
	resp.ItemMedia = res
	return c.JSON(resp)
}
