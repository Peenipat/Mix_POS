package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"	
	// "myapp/modules/core/models"
	"log"
	// bookingModels "myapp/modules/barberbooking/models"
	"myapp/seeds"
	"os"

)
// const (
// 	host = "localhost"
// 	port = 5432
// 	user = "myuser"
// 	password = "mypassword"
// 	dbname = "mydatabase"
// )

// var DB *gorm.DB

// func ConnectDB(){
// 	dsn := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
// 	db ,err  := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil{
// 		panic("failed to connect to database")
// 	} 

var DB *gorm.DB

func ConnectDB() {
    // 1) อ่าน DATABASE_URL จาก env ถ้ามี
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("Missing DATABASE_URL")
    }

    log.Printf("→ DSN: %s\n", dsn)

	var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }

	// if err := DB.AutoMigrate(
    //     &coreModels.Tenant{},
    //     &coreModels.Role{},
    //     &coreModels.Branch{},
    //     &coreModels.TenantUser{},
    //     &coreModels.User{},
    //     &coreModels.SystemLog{},
    //     &bookingModels.Service{},
    //     &bookingModels.WorkingHour{},
    //     &bookingModels.Unavailability{},
    //     &bookingModels.Barber{},
    //     &bookingModels.Appointment{},
    // ); err != nil {
    //     log.Fatalf("failed to auto migrate models: %v", err)
    // }
    // 0) Seed Roles
    if err := seeds.SeedModules(DB); err != nil {
        log.Fatalf("failed to seed roles: %v", err)
    }
     // 1) Seed Roles
     if err := seeds.SeedRoles(DB); err != nil {
        log.Fatalf("failed to seed roles: %v", err)
    }
    // 2) Seed Tenants
    if err := seeds.SeedTenants(DB); err != nil {
        log.Fatalf("failed to seed tenants: %v", err)
    }
    // 3) Seed Branch
    if err := seeds.SeedBranches(DB); err != nil {
        log.Fatalf("seed branches failed: %v", err)
    }
    // 4) Seed User
    if err := seeds.SeedUsers(DB); err != nil {
        log.Fatalf("seed users failed: %v", err)
    }
    // 5) Seed TenantUsers
    if err := seeds.SeedTenantUsers(DB); err != nil {
        log.Fatalf("seed tenant_users failed: %v", err)
    }
    // 6) Seed Services
    if err := seeds.SeedServices(DB); err != nil {
        log.Fatalf("seed services failed: %v", err)
    }
    // 7) Seed WorkingHours
    if err := seeds.SeedWorkingHours(DB); err != nil {
        log.Fatalf("seed working hours failed: %v", err)
    }
     // 8) Seed Unavailabilities
    if err := seeds.SeedUnavailabilities(DB); err != nil {
        log.Fatalf("seed unavailabilities failed: %v", err)
    }
      // 9) Seed Barbers
    if err := seeds.SeedBarbers(DB); err != nil {
        log.Fatalf("seed barbers failed: %v", err)
    }
      // 10) Seed Appointments
    if err := seeds.SeedAppointments(DB); err != nil {
        log.Fatalf("seed appointments failed: %v", err)
    }
	

	
}