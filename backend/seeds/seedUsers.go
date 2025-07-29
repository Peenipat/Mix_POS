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

    // โหลด Default Branch 
    var branch coreModels.Branch
    if err := db.Where("name = ?", "Branch 1").First(&branch).Error; err != nil {
        return err
    }

    var branch2 coreModels.Branch
    if err := db.Where("name = ?", "Branch 2").First(&branch2).Error; err != nil {
        return err
    }

    // กำหนดรายการผู้ใช้ตั้งต้น
    type seed struct {
        Username string
        Email    string
        Password string
        RoleID   uint
        BranchID *uint
        Img_path string
        Img_name string
    }
    data := []seed{
        // SaaS SuperAdmin (ไม่ผูกสาขา)
        {"saas_admin", "saas_admin@gmail.com", "12345678", roleSA.ID, nil,"barbers","barber1.jpg"},
        // Tenant Admin (ดูแลหลายสาขา)
        {"tenant_admin", "tenant_admin@gmail.com", "12345678", roleTA.ID, &branch.ID,"barbers","barber3.jpg"},
        // Branch Admin (เฉพาะสาขา)
        {"branch_admin", "branch_admin@gmail.com", "12345678", roleBA.ID, &branch.ID,"barbers","barber2.jpg"},
        {"branch_admin2", "branch2_admin@gmail.com", "12345678", roleBA.ID, &branch2.ID,"barbers","barber2.jpg"},
        // Assistant Manager
        {"assistant_mgr", "assistant@gmail.com", "12345678", roleAM.ID, &branch.ID,"barbers","barber4.jpg"},
        // Staff
        {"staff_user", "staff@gmail.com", "12345678", roleST.ID, &branch.ID,"barbers","barber2.jpg"},
        // End-customer / general user
        {"generic_user", "user@gmail.com", "12345678", roleUS.ID, &branch.ID,"barbers","barber1.jpg"},
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
            Img_path: u.Img_path,
            Img_name: u.Img_name,
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
