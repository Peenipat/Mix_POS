package seeds

import (
    "time"

    "gorm.io/gorm"
    bookingModels "myapp/modules/barberbooking/models"
)

func SeedServices(db *gorm.DB) error {
    items := []bookingModels.Service{
        {Name: "Haircut",   Description:"ตัดผม", Duration: 30, Price: 200, Img_path:"service",Img_name:"service1.jpg",TenantID:	1,BranchID:	1},
        {Name: "Shampoo",   Description:"ตัดผม", Duration: 15, Price: 100, Img_path:"service",Img_name:"service2.jpg",TenantID:	1,BranchID:	1},
        {Name: "Beard Trim",Description:"ตัดผม", Duration: 45, Price: 150, Img_path:"service",Img_name:"service3.jpg",TenantID:	1,BranchID:	1},
        {Name: "Beard Trim2",Description:"ตัดผม", Duration: 20, Price: 150, Img_path:"service",Img_name:"service4.jpg",TenantID: 1,BranchID: 2},
    }

    now := time.Now()
    for _, svc := range items {
        record := bookingModels.Service{Name: svc.Name}
        attrs  := bookingModels.Service{
            Description:svc.Description,
            Duration:  svc.Duration,  
            Price:     svc.Price,
            Img_path: svc.Img_path,
            Img_name: svc.Img_name,
            TenantID: svc.TenantID,
            BranchID: svc.BranchID,
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
