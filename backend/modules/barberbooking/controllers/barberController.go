package barberBookingController

import (
	"log"
	helperFunc "myapp/modules/barberbooking"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BarberController struct {
	BarberService barberBookingPort.IBarber
}

func NewBarberController(svc barberBookingPort.IBarber) *BarberController {
	return &BarberController{
		BarberService: svc,
	}
}



// GetBarberByID godoc
// @Summary      ดึงข้อมูลช่างตัดผมตาม ID
// @Description  คืนข้อมูล Barber ตามรหัสที่ระบุ
// @Tags         Barber
// @Accept       json
// @Produce      json
// @Param        barber_id  path      uint                             true  "รหัส Barber"
// @Success      200        {object}  map[string]interface{}          "คืนค่า status success, message และข้อมูลช่างตัดผมใน key `data`"
// @Failure      400        {object}  map[string]string               "Invalid barber_id"
// @Failure      404        {object}  map[string]string               "Barber not found"
// @Failure      500        {object}  map[string]string               "Failed to fetch barber"
// @Router       /barbers/{barber_id} [get]
func (ctrl *BarberController) GetBarberByID(c *fiber.Ctx) error {

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
// ListBarbersByBranch godoc
// @Summary      ดึงรายชื่อช่างตัดผมตามสาขา
// @Description  คืนรายการ Barber ทั้งหมดของสาขาที่ระบุ หากไม่พบจะคืน 404 พร้อม data ว่าง
// @Tags         Barber
// @Accept       json
// @Produce      json
// @Param        branch_id  path      uint                                   true  "รหัส Branch"
// @Success      200        {object}  map[string]interface{}               "คืนค่า status success, message และ array ของ Barber ใน key `data`"
// @Failure      400        {object}  map[string]string                    "invalid branch_id"
// @Failure      404        {object}  map[string]interface{}               "no barbers found for this branch (data จะเป็น array ว่าง)"
// @Failure      500        {object}  map[string]string                    "failed to fetch barbers"
// @Router       /branches/{branch_id}/barbers [get]
// @Security     ApiKeyAuth
func (ctrl *BarberController) ListBarbersByBranch(c *fiber.Ctx) error {
    branchID, err := helperFunc.ParseUintParam(c, "branch_id")
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid branch_id",
            "error":   err.Error(),
        })
    }

    // 2. Call the service
    barbers, err := ctrl.BarberService.ListBarbersByBranch(c.Context(), &branchID)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "failed to fetch barbers",
            "error":   err.Error(),
        })
    }

    // 3. If no barbers found, 404
    if len(barbers) == 0 {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{
            "status":  "error",
            "message": "no barbers found for this branch",
            "data":    []barberBookingPort.BarberWithUser{},
        })
    }

    // 4. Return the list
    return c.Status(http.StatusOK).JSON(fiber.Map{
        "status":  "success",
        "message": "barber list retrieved",
        "data":    barbers,
    })
}

var RolesCanManageBarber = []coreModels.RoleName{
	coreModels.RoleNameTenantAdmin,
	coreModels.RoleNameTenant,
	coreModels.RoleNameBranchAdmin,
}


// CreateBarber godoc
// @Summary      สร้างช่างตัดผมใหม่
// @Description  เพิ่ม Barber ใหม่ลงในระบบ (ต้องมีสิทธิ์ TenantAdmin, Tenant หรือ BranchAdmin)
// @Tags         Barber
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint                             true  "รหัส tenant"
// @Param        branch_id  path      uint                             true  "รหัส branch"
// @Param        body  body      barberBookingPort.CreateBarberInput  true  "Payload สำหรับสร้าง Barber (UserID, BranchID, ชื่อ-นามสกุล ฯลฯ)"
// @Success      201   {object}  map[string]string          "คืนค่า status success และข้อความยืนยันการสร้าง"
// @Failure      400   {object}  map[string]string          "Invalid request body"
// @Failure      403   {object}  map[string]string          "Permission denied"
// @Failure      500   {object}  map[string]string          "Failed to create barber"
// @Router       /tenants/{tenant_id}/branches/{branch_id}/create-barber [post]
// @Security     ApiKeyAuth
func (ctrl *BarberController) CreateBarber(c *fiber.Ctx) error {
    // 1) ตรวจ permission
    roleStr, ok := c.Locals("role").(string)
    log.Println("Role :", roleStr)

    if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageBarber) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "status":  "error",
            "message": "Permission denied",
        })
    }


    tenantIDParam := c.Params("tenant_id")
    tenantIDUint64, err := strconv.ParseUint(tenantIDParam, 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant_id",
        })
    }
    tenantID := uint(tenantIDUint64)
    // 3) ดึง branch_id จาก path param
    branchIDParam := c.Params("branch_id")
    branchIDUint64, err := strconv.ParseUint(branchIDParam, 10, 64)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid branch_id",
        })
    }
    branchID := uint(branchIDUint64)

    // 4) Parse body ให้โครงสร้างเฉพาะ user_id กับ phone_number
    var body barberBookingPort.CreateBarberInput
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid request body",
        })
    }
    if body.UserID == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "user_id is required",
        })
    }
    if strings.TrimSpace(body.PhoneNumber) == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "phone_number is required",
        })
    }


    // 5) สร้าง payload ของ Barber
    payload := &barberBookingModels.Barber{
        TenantID:    tenantID,
        BranchID:    branchID,
        UserID:      body.UserID,
        Description: body.Description,
        PhoneNumber: strings.TrimSpace(body.PhoneNumber),
    }

    // 6) เรียก service สร้าง
    if err := ctrl.BarberService.CreateBarber(c.Context(), payload); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to create barber",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "status":  "success",
        "message": "Barber created",
    })
}

// UpdateBarber godoc
// @Summary      แก้ไขข้อมูลช่างตัดผม
// @Description  อัปเดตข้อมูล Barber ตามรหัสที่ระบุ (ต้องมีสิทธิ์ TenantAdmin, Tenant หรือ BranchAdmin)
// @Tags         Barber
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint                              true  "รหัส Tenant"
// @Param        barber_id  path      uint                              true  "รหัส Barber"
// @Param        body       body      barberBookingPort.UpdateBarberRequest  true  "Payload สำหรับอัปเดต Barber (BranchID, UserID, PhoneNumber, Username, Email)"
// @Success      200        {object}  map[string]interface{}            "คืนค่า status success และข้อมูล Barber ใน key `data`"
// @Failure      400        {object}  map[string]string                "Invalid tenant_id/barber_id หรือ Invalid request body"
// @Failure      403        {object}  map[string]string                "Permission denied"
// @Failure      404        {object}  map[string]string                "Barber not found"
// @Failure      500        {object}  map[string]string                "Failed to update Barber"
// @Router       /tenants/{tenant_id}/barbers/{barber_id}/update-barber [put]
// @Security     ApiKeyAuth
func (ctrl *BarberController) UpdateBarber(c *fiber.Ctx) error {
    // 1. เช็คสิทธิ์
    roleStr, ok := c.Locals("role").(string)
    if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageBarber) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "status":  "error",
            "message": "Permission denied",
        })
    }

    // 2. ดึง barber_id จาก path parameter
    barberID, err := helperFunc.ParseUintParam(c, "barber_id")
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid barber_id",
            "error":   err.Error(),
        })
    }

    // 3. แปลง body เป็น DTO
    var payload barberBookingPort.UpdateBarberRequest
    if err := c.BodyParser(&payload); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid request body",
            "error":   err.Error(),
        })
    }

    // 4. ดึงข้อมูล barber ปัจจุบัน
    existingBarber, err := ctrl.BarberService.GetBarberByID(c.Context(), barberID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.Status(http.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "Barber not found",
            })
        }
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to fetch barber",
            "error":   err.Error(),
        })
    }
    if existingBarber == nil {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{
            "status":  "error",
            "message": "Barber not found",
        })
    }

    // 5. เซ็ตค่าที่จะแก้ไขลงใน existingBarber
    existingBarber.BranchID = payload.BranchID
    existingBarber.PhoneNumber = payload.PhoneNumber
    existingBarber.Description = payload.Description

    // 6. เรียก Service ให้ทำการอัปเดตทั้งใน table barbers และอัปเดตชื่อ-อีเมลในตาราง users
    updatedBarber, err := ctrl.BarberService.UpdateBarber(
        c.Context(),
        barberID,
        existingBarber,
        payload.Username,
        payload.Email,
    )
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to update barber",
            "error":   err.Error(),
        })
    }

    // 7. ส่งผลลัพธ์กลับ
    return c.Status(http.StatusOK).JSON(fiber.Map{
        "status":  "success",
        "message": "Barber updated",
        "data":    updatedBarber,
    })
}

// DeleteBarber godoc
// @Summary      ลบช่างตัดผม
// @Description  ลบ Barber ตามรหัสที่ระบุ (ต้องมีสิทธิ์ TenantAdmin, Tenant หรือ BranchAdmin)
// @Tags         Barber
// @Accept       json
// @Produce      json
// @Param        barber_id  path      uint                   true  "รหัส Barber"
// @Success      200        {object}  map[string]string      "คืนค่า status success และข้อความยืนยันการลบ"
// @Failure      400        {object}  map[string]string      "Invalid barber_id"
// @Failure      403        {object}  map[string]string      "Permission denied"
// @Failure      500        {object}  map[string]string      "Failed to delete barber"
// @Router       /barbers/{barber_id} [delete]
// @Security     ApiKeyAuth
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

// GetBarberByUser godoc
// @Summary      ดึงข้อมูลช่างตัดผมโดย User ID
// @Description  คืนข้อมูล Barber ที่เชื่อมโยงกับ User ที่ระบุ (ต้องมีสิทธิ์ TenantAdmin, Tenant หรือ BranchAdmin)
// @Tags         Barber
// @Accept       json
// @Produce      json
// @Param        user_id  path      uint                             true  "รหัส User"
// @Success      200      {object}  map[string]interface{}          "คืนค่า status success, message และข้อมูล Barber ใน key `data`"
// @Failure      400      {object}  map[string]string               "Invalid user_id"
// @Failure      403      {object}  map[string]string               "Permission denied"
// @Failure      404      {object}  map[string]string               "Barber not found"
// @Failure      500      {object}  map[string]string               "Failed to fetch barber"
// @Router       /users/{user_id}/barber [get]
// @Security     ApiKeyAuth
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

var RoleCanGetBarberByTenant = []coreModels.RoleName{
	coreModels.RoleNameTenantAdmin,
	coreModels.RoleNameTenant,
}
// ListBarbersByTenant godoc
// @Summary      ดึงรายชื่อช่างตาม Tenant
// @Description  คืนรายการ Barber ทั้งหมดของ Tenant ที่ระบุ (ต้องมีสิทธิ์ TenantAdmin หรือ Tenant)
// @Tags         Barber
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint                              true   "รหัส Tenant"
// @Success      200        {object}  map[string]interface{}           "คืนค่า status, message และ array ของ Barber ใน key `data`"
// @Failure      400        {object}  map[string]string                "Invalid tenant_id"
// @Failure      403        {object}  map[string]string                "Permission denied"
// @Failure      404        {object}  map[string]string                "List Barber not found"
// @Failure      500        {object}  map[string]string                "Failed to fetch List Barber"
// @Router       /tenants/{tenant_id}/barbers [get]
// @Security     ApiKeyAuth
func (ctrl *BarberController) ListBarbersByTenant(c *fiber.Ctx) error{
	roleStr,ok := c.Locals("role").(string)
    log.Println("Role :", roleStr)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr,RoleCanGetBarberByTenant){
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

	if listBarber == nil{
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":"error",
			"message":"List list Barber not found",
		})
	}

	if len(listBarber) == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "List Barber not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Barbers retrieved",
		"data":    listBarber,
	})


}

// ListUserNotBarber godoc
// @Summary         ดึงข้อมูล User ไม่ได้เป็น barber ในสาขานั้น ๆ 
// @Description     คืนรายการ user ที่ไม่ได้เป็น Barber ทั้งหมดของสาขานั้น
// @Tags            Barber
// @Accept          json
// @Produce         json
// @Param           tenant_id  path      uint                              true   "รหัส Tenant"
// @Param           branch_id  path      uint                              true   "รหัส Branch"
// @Success      200        {object}  map[string]interface{}           "คืนค่า status, message และ array ของ Barber ใน key `data`"
// @Failure      400        {object}  map[string]string                "Invalid tenant_id"
// @Failure      403        {object}  map[string]string                "Permission denied"
// @Failure      404        {object}  map[string]string                "List Barber not found"
// @Failure      500        {object}  map[string]string                "Failed to fetch List Barber"
// @Router          /tenants/{tenant_id}/branches/{branch_id}/users [get]
// @Security        ApiKeyAuth // ถ้าใช้ token
func (ctrl *BarberController) ListUserNotBarber(c *fiber.Ctx) error {
    // 1. Parse the branch_id URL param
    branchID, err := helperFunc.ParseUintParam(c, "branch_id")
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid branch_id",
            "error":   err.Error(),
        })
    }

    // 2. Call the service
    usersNotBarber, err := ctrl.BarberService.ListUserNotBarber(c.Context(), &branchID)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "failed to fetch barbers",
            "error":   err.Error(),
        })
    }


    // 4. Return the list
    return c.Status(http.StatusOK).JSON(fiber.Map{
        "status":  "success",
        "message": "barber list retrieved",
        "data":    usersNotBarber,
    })
}
