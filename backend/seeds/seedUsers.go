package seeds

import (
    "time"

    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    coreModels "myapp/modules/core/models"
)

// SeedUsers สร้างผู้ใช้ตั้งต้นให้ครบทุกรายการบทบาท
func SeedUsers(db *gorm.DB) error {
    // โหลด Role records
    var (
        roleSA coreModels.Role
        roleTA coreModels.Role
        roleBA coreModels.Role
        roleAM coreModels.Role
        roleST coreModels.Role
        roleUS coreModels.Role
    )
    roles := map[string]*coreModels.Role{
        string(coreModels.RoleNameSaaSSuperAdmin):   &roleSA,
        string(coreModels.RoleNameTenantAdmin):      &roleTA,
        string(coreModels.RoleNameBranchAdmin):      &roleBA,
        string(coreModels.RoleNameAssistantManager): &roleAM,
        string(coreModels.RoleNameStaff):            &roleST,
        string(coreModels.RoleNameUser):             &roleUS,
    }
    for name, ptr := range roles {
        if err := db.Where("name = ?", name).First(ptr).Error; err != nil {
            return err
        }
    }

    // โหลด Default Branch (ใช้สำหรับทุกบทบาทที่ผูกสาขา)
    var branch coreModels.Branch
    if err := db.Where("name = ?", "Default Branch").First(&branch).Error; err != nil {
        return err
    }

    // กำหนดรายการผู้ใช้ตั้งต้น
    type seed struct {
        Username string
        Email    string
        Password string
        RoleID   uint
        BranchID *uint
    }
    data := []seed{
        // SaaS SuperAdmin (ไม่ผูกสาขา)
        {"saas_admin", "saas_admin@yourdomain.com", "supersecret", roleSA.ID, nil},
        // Tenant Admin (ดูแลหลายสาขา)
        {"tenant_admin", "tenant_admin@default.example.com", "tenantsecret", roleTA.ID, nil},
        // Branch Admin (เฉพาะสาขา)
        {"branch_admin", "branch_admin@default.example.com", "branchsecret", roleBA.ID, &branch.ID},
        // Assistant Manager
        {"assistant_mgr", "assistant@default.example.com", "assistsecret", roleAM.ID, &branch.ID},
        // Staff
        {"staff_user", "staff@default.example.com", "staffsecret", roleST.ID, &branch.ID},
        // End-customer / general user
        {"generic_user", "user@default.example.com", "usersecret", roleUS.ID, nil},
    }

    now := time.Now()
    for _, u := range data {
        // hash password
        hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }

        record := coreModels.User{Email: u.Email}
        attrs := coreModels.User{
            Username:  u.Username,
            Password:  string(hashed),
            RoleID:    u.RoleID,
            BranchID:  u.BranchID,
            UpdatedAt: now,
        }
        if err := db.Where(record).
            Assign(attrs).
            FirstOrCreate(&record).Error; err != nil {
            return err
        }
    }

    return nil
}
