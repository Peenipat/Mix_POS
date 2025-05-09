package barberbooking
import (
		"github.com/gofiber/fiber/v2"
		"strconv"
		"fmt"
		coreModels "myapp/modules/core/models"
)

func ParseUintParam(c *fiber.Ctx, paramName string) (uint, error) {
	val := c.Params(paramName)
	id, err := strconv.Atoi(val)
	if err != nil || id < 0 {
		return 0, fmt.Errorf("invalid %s", paramName)
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