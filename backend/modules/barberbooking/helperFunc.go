package barberbooking
import (
		"github.com/gofiber/fiber/v2"
		"strconv"
		"fmt"
		"database/sql/driver"
		"time"
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


type TimeOnly struct {
	time.Time
}

func (t *TimeOnly) UnmarshalJSON(b []byte) error {
	parsed, err := time.Parse(`"15:04"`, string(b))
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func (t TimeOnly) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Format("15:04"))), nil
}

func (t *TimeOnly) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	case string:
		parsed, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into TimeOnly", value)
	}
}

func (t TimeOnly) Value() (driver.Value, error) {
	return t.Format("15:04:05"), nil
}