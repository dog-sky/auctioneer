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

	realmID := h.BlizzClient.Cache.GetRealmID(queryParams.RealmName)
	if realmID == 0 {
		return fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Sprintf("Realm %s not found", queryParams.RealmName),
		)
	}

	res := new(ResponseV1)
	res.Success = true

	// TODO
	res.Message = "DEBUG"
	// Получить в структуру все айди предметов, которые соответсвуют поиску (или пропустить этот шаг если не передано значение)
	// Потом сделать запрос в аук по серверу
	// получить все медиа айтимов, которые пересекаются по айди (или вообще всех?) Нужны ли медиа вообще, или можно пока выдавать результат текстом?
	// отдать в ответе только то, что пересекается с запрошеным предметом (или всё, если запрос пустой)

	// namespace := fmt.Sprintf("dynamic-%s", queryParams.region)

	return c.JSON(res)
}
