package barberbooking
import (
		"github.com/gofiber/fiber/v2"
		"strconv"
		"fmt"
		coreModels "myapp/modules/core/models"
)


func ParseUintParam(c *fiber.Ctx, name string) (uint, error) {
	param := c.Params(name)
	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid uint param: %s", param)
	}
	return uint(id), nil
}

func IsAuthorizedRole(role string, allowed []coreModels.RoleName) bool {
	for _, r := range allowed {
		if string(r) == role {
			return true
		}
	}
	return false
}