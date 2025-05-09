package barberBookingController

import (
	helperFunc "myapp/modules/barberbooking"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type BarberController struct {
	BarberService barberBookingPort.IBarber
}

func NewBarberController(svc barberBookingPort.IBarber) *BarberController {
	return &BarberController{
		BarberService: svc,
	}
}

var RolesCanManageBarber = []coreModels.RoleName{
	coreModels.RoleNameSaaSSuperAdmin,
	coreModels.RoleNameTenant,
	coreModels.RoleNameBranchAdmin,
}

func (ctrl *BarberController) CreateBarber(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageBarber) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	var payload barberBookingModels.Barber
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if err := ctrl.BarberService.CreateBarber(c.Context(), &payload); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create barber",
			"error":   "Can't create barber",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Customer created",
	})

}

func (ctrl *BarberController) GetBarberByID(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageBarber) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	barberID, err := helperFunc.ParseUintParam(c, "barber_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	barber, err := ctrl.BarberService.GetBarberByID(c.Context(), barberID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch barber",
			"error":   "Internal server error",
		})
	}

	if barber == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":"error",
			"message":"Barber not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":"success",
		"message":"Barber retrieved",
		"data":barber,
	})

}

func (ctrl *BarberController) ListBarbersByBranch(c *fiber.Ctx) error{
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr,RolesCanManageBarber){
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":"error",
			"message":"Permission denied",
		})
	}

	branchID, err := helperFunc.ParseUintParam(c,"branch_id")
	if err != nil{
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	listBarber,err := ctrl.BarberService.ListBarbersByBranch(c.Context(),&branchID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":"error",
			"message":"Failed to fetch List Barber",
			"error":"Internal server error",
		})
	}

	if listBarber == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":"error",
			"message":"List Barber not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":"success",
		"message":"List Barber retrieved",
		"data": listBarber,
	})
}

func (ctrl *BarberController) UpdateBarber(c *fiber.Ctx) error{
	roleStr,ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr,RolesCanManageBarber){
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":"error",
			"message":"Permission denied",
		})
	}

	barberID, err := helperFunc.ParseUintParam(c, "barber_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var payload barberBookingModels.Barber
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":"error",
			"message":"Invalid request body",
		})
	}

	existingBarber, err := ctrl.BarberService.GetBarberByID(c.Context(),barberID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Barber not found",
		})
	}

	existingBarber.BranchID = payload.BranchID
	updateBarber, err := ctrl.BarberService.UpdateBarber(c.Context(),barberID,existingBarber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update Barber",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":"success",
		"message":"Barber Updated",
		"data":updateBarber,
	})

}

func (ctrl *BarberController) DeleteBarber(c *fiber.Ctx) error{
	roleStr,ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr,RolesCanManageBarber){
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":"error",
			"message":"Permission denied",
		})
	}

	barberID, err := helperFunc.ParseUintParam(c, "barber_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := ctrl.BarberService.DeleteBarber(c.Context(),barberID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":"error",
			"message":"Failed to delete barber",
			"error":err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":"success",
		"message":"Barber delete successfully",
	})
}

func (ctrl *BarberController) GetBarberByUser(c *fiber.Ctx) error{
	roleStr,ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr,RolesCanManageBarber){
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":"error",
			"message":"Permission denied",
		})
	}

	userID, err := helperFunc.ParseUintParam(c, "user_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	barber,err := ctrl.BarberService.GetBarberByUser(c.Context(),userID)
	if err != nil{
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":"error",
			"message":"Feiled to fetch barber",
			"error":"Internal server error",
		})
	}

	if barber == nil{
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":"error",
			"message":"Barber not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":"success",
		"message":"Barber retrieved",
		"data": barber,
	})	
}

func (ctrl *BarberController) ListBarbersByTenant(c *fiber.Ctx) error{
	roleStr,ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr,RolesCanManageBarber){
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":"error",
			"message":"Permission denied",
		})
	}

	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	listBarber, err := ctrl.BarberService.ListBarbersByTenant(c.Context(),tenantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":"error",
			"message":"Failed to fetch List Barber",
			"error":"Internal server error",
		})
	}

	if listBarber == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":"error",
			"message":"List list Barber not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":"success",
		"message":"Customer retrieved",
		"data": listBarber,
	})


}