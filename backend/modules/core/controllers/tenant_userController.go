package Core_controllers

import (
	"errors"
	"strconv"

	corePort "myapp/modules/core/port"
	coreServices "myapp/modules/core/services"
    coreModels "myapp/modules/core/models"
    helperFunc "myapp/modules/core"

	"github.com/gofiber/fiber/v2"
)

// TenantUserController handles M2M endpoints between tenants and users.
type TenantUserController struct {
    Service corePort.ITenantUser
}

// NewTenantUserController constructs a new controller.
func NewTenantUserController(svc corePort.ITenantUser) *TenantUserController {
    return &TenantUserController{Service: svc}
}

var (
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrUserAlreadyAssigned = errors.New("user already assigned to tenant")
	ErrUserNotAssigned = errors.New("user not assigned to tenant")
    ErrNoTenantsAssigned  = errors.New("no tenants assigned to user")
    ErrUserNotFound       = errors.New("user not found")
)

// AddUserToTenant godoc
// @Summary      เพิ่มผู้ใช้งานให้กับ Tenant
// @Description  กำหนดให้ User ที่ระบุด้วย `user_id` เข้าใช้งาน Tenant ที่ระบุด้วย `tenant_id`
// @Tags         TenantUser
// @Produce      json
// @Param        tenant_id  path      int  true  "รหัส Tenant"
// @Param        user_id    path      int  true  "รหัส User"
// @Success      201        {object}  map[string]string  "คืนค่า status และข้อความยืนยันการเพิ่มผู้ใช้"
// @Failure      400        {object}  map[string]string  "Invalid tenant ID หรือ invalid user ID"
// @Failure      404        {object}  map[string]string  "Tenant not found หรือ User not found"
// @Failure      409        {object}  map[string]string  "User already assigned to this tenant"
// @Failure      500        {object}  map[string]string  "เกิดข้อผิดพลาดระหว่างการเพิ่มผู้ใช้งานให้ Tenant"
// @Router       /core/tenant-user/tenants/:tenant_id/users/:user_id [post]
// @Security     ApiKeyAuth
func (ctrl *TenantUserController) AddUserToTenant(c *fiber.Ctx) error {
    tidParam := c.Params("tenant_id")
    tid64, err := strconv.ParseUint(tidParam, 10, 64)
    if err != nil || tid64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant ID",
        })
    }

    // 3. Parse user_id
    uidParam := c.Params("user_id")
    uid64, err := strconv.ParseUint(uidParam, 10, 64)
    if err != nil || uid64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid user ID",
        })
    }

    // 4. Call service
    err = ctrl.Service.AddUserToTenant(c.Context(), uint(tid64), uint(uid64))
    if err != nil {
        switch {
        case errors.Is(err, coreServices.ErrInvalidTenantID),
             errors.Is(err, coreServices.ErrInvalidUserID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrTenantNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "Tenant not found",
            })
        case errors.Is(err, coreServices.ErrUserNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "User not found",
            })
        case errors.Is(err, coreServices.ErrUserAlreadyAssigned):
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "status":  "error",
                "message": "User already assigned to this tenant",
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to assign user to tenant",
            })
        }
    }

    // 5. Success
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "status":  "success",
        "message": "User assigned to tenant",
    })
}


// RemoveUserFromTenant godoc
// @Summary      นำผู้ใช้ออกจาก Tenant
// @Description  ลบการผูก User ที่ระบุด้วย `user_id` ออกจาก Tenant ที่ระบุด้วย `tenant_id`
// @Tags         TenantUser
// @Produce      json
// @Param        tenant_id  path      int  true  "รหัส Tenant"
// @Param        user_id    path      int  true  "รหัส User"
// @Success      200        {object}  map[string]string  "คืนค่า status และข้อความยืนยันการลบผู้ใช้"
// @Failure      400        {object}  map[string]string  "Invalid tenant ID หรือ Invalid user ID"
// @Failure      404        {object}  map[string]string  "Tenant not found หรือ User not found"
// @Failure      409        {object}  map[string]string  "User is not assigned to this tenant"
// @Failure      500        {object}  map[string]string  "เกิดข้อผิดพลาดระหว่างการลบผู้ใช้จาก Tenant"
// @Router       /core/tenant-user/tenants/:tenant_id/users/:user_id [delete]
// @Security     ApiKeyAuth
func (ctrl *TenantUserController) RemoveUserFromTenant(c *fiber.Ctx) error {
    // 1. Parse tenant_id
    tidParam := c.Params("tenant_id")
    tid64, err := strconv.ParseUint(tidParam, 10, 64)
    if err != nil || tid64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant ID",
        })
    }
    tenantID := uint(tid64)

    // 2. Parse user_id
    uidParam := c.Params("user_id")
    uid64, err := strconv.ParseUint(uidParam, 10, 64)
    if err != nil || uid64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid user ID",
        })
    }
    userID := uint(uid64)

    // 3. Call service
    err = ctrl.Service.RemoveUserFromTenant(c.Context(), tenantID, userID)
    if err != nil {
        switch {
        case errors.Is(err, coreServices.ErrInvalidTenantID),
             errors.Is(err, coreServices.ErrInvalidUserID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrTenantNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "Tenant not found",
            })
        case errors.Is(err, coreServices.ErrUserNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "User not found",
            })
        case errors.Is(err, coreServices.ErrUserNotAssigned):
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "status":  "error",
                "message": "User is not assigned to this tenant",
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to remove user from tenant",
            })
        }
    }

    // 4. Success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":  "success",
        "message": "User removed from tenant",
    })
}


// ListTenantsByUser godoc
// @Summary      ดึงรายการ Tenant ของ User
// @Description  ดึงรายการ Tenant ที่ User ระบุเข้าใช้งานได้ (user_id มาจาก path)
// @Tags         TenantUser
// @Produce      json
// @Param        user_id  path      int  true  "รหัส User"
// @Success      200      {object}  map[string]interface{}  "คืนค่า status และ array ของ Tenant ใน key `data`"
// @Failure      400      {object}  map[string]string       "Invalid user ID"
// @Failure      404      {object}  map[string]string       "User not found"
// @Failure      500      {object}  map[string]string       "เกิดข้อผิดพลาดระหว่างดึงรายการ Tenant"
// @Router       /core/tenant-user/user/:user_id [get]
// @Security     ApiKeyAuth
func (ctrl *TenantUserController) ListTenantsByUser(c *fiber.Ctx) error {
    // 1. Parse and validate user_id
    uidParam := c.Params("user_id")
    uid64, err := strconv.ParseUint(uidParam, 10, 64)
    if err != nil || uid64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid user ID",
        })
    }
    userID := uint(uid64)

    // 2. Call service
    tenants, err := ctrl.Service.ListTenantsByUser(c.Context(), userID)
    if err != nil {
        // 3. Handle service errors
        switch {
        case errors.Is(err, coreServices.ErrInvalidUserID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrUserNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "User not found",
            })
        case errors.Is(err, coreServices.ErrNoTenantsAssigned):
            // เปลี่ยนจาก 404+error เป็น 200+success พร้อม data ว่าง
            return c.JSON(fiber.Map{
                "status": "success",
                "data":   []coreModels.Tenant{},
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to fetch tenants",
            })
        }
    }

    // 4. Success
    return c.JSON(fiber.Map{
        "status": "success",
        "data":   tenants,
    })
}


var RolesCanManageTenantUsers = []coreModels.RoleName{
	coreModels.RoleNameTenantAdmin,
    coreModels.RoleNameBranchAdmin,
	coreModels.RoleNameTenant,
}
// ListUsersForTenant godoc
// @Summary      List users under a Tenant
// @Description  Returns all users associated with the specified tenant ID.
// @Tags         TenantUsers
// @Accept       json
// @Produce      json
// @Param        tenant_id   path      uint     true   "Tenant ID"
// @Success      200         {object}  map[string][]UserResponse   "status=success, data=array of users"
// @Failure      400         {object}  map[string]string           "Invalid tenant_id"
// @Failure      403         {object}  map[string]string           "Permission denied"
// @Failure      500         {object}  map[string]string           "Failed to list users"
// @Router       /core/tenant-user/tenants/:tenant_id [get]
// @Security     ApiKeyAuth
func (ctrl *TenantUserController) ListUsersForTenant(c *fiber.Ctx) error {
    // ตรวจ permission …
    roleStr, ok := c.Locals("role").(string)
    if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageTenantUsers) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "status":  "error",
            "message": "Permission denied",
        })
    }

    // Parse tenant_id จาก path param
    tenantIDParam := c.Params("tenant_id")
    tid, err := strconv.ParseUint(tenantIDParam, 10, 64)
    if err != nil || tid == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant_id",
        })
    }
    tenantID := uint(tid)

    // เรียก service มาดึง users
    users, err := ctrl.Service.ListUsersByTenant(c.Context(), tenantID)
    if err != nil {
        // ส่งกลับ HTTP 500 พร้อมข้อความ error
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to list users",
            "error":   err.Error(),
        })
    }

    // ถ้า len(users) == 0 ก็ยัง respond 200 และ data เป็น slice ว่าง
    var resp []map[string]interface{}
    for _, u := range users {
        resp = append(resp, map[string]interface{}{
            "id":       u.ID,
            "username": u.Username,
            "email":    u.Email,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   resp,
    })
}
