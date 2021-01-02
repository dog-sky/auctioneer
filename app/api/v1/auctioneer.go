package v1

import (
	"auctioneer/app/blizz"
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
)

func (h *V1Handler) SearchItemData(c *fiber.Ctx) error {
	params := new(searchQueryParams)
	if err := c.QueryParser(params); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Sprintf("Invalid query params. Err: %v", err),
		)
	}

	if err := checkQueryParams(params); err != nil {
		return err
	}

	realmID := h.BlizzClient.GetRealmID(params.RealmName)
	if realmID == 0 {
		return fiber.NewError(
			fiber.StatusNotFound,
			fmt.Sprintf("Realm %s not found", params.RealmName),
		)
	}

	res := new(ResponseV1)
	res.Success = true

	itemData := new(blizz.ItemResult)
	if params.ItemName != "" {
		searchResult, err := h.BlizzClient.SearchItem(params.ItemName, params.Region)
		if err != nil {
			return fiber.NewError(
				fiber.StatusBadRequest,
				err.Error(),
			)
		}
		if len(searchResult.Results) == 0 {
			return fiber.NewError(
				fiber.StatusNotFound,
				fmt.Sprintf("Item %s not found", params.ItemName),
			)
		}
		itemData = searchResult
	}

	data, err := h.BlizzClient.GetAuctionData(realmID, params.Region)
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	if len(itemData.Results) == 0 {
		res.Result = data
		return c.JSON(res)
	}

	for _, AucItem := range data {
		for _, item := range itemData.Results {
			if AucItem.Item.ID == item.Data.ID {
				res.Result = append(res.Result, AucItem)
			}
		}
	}

	// TODO В ответ в модели нужно ещё выдавать имя предмета и его качество.
	// возможно нам нужна функциональность только поиска по предмету и не делать функциональность выгрузки всего аука
	// Если ничего не нашел -- вернуть пустой резулт. Сейчас не возмращает ничего
	return c.JSON(res)
}
