// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"os"
	"time"

	// "net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	// "github.com/gofiber/fiber/v2/middleware/filesystem"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"myapp/database"
	_ "myapp/docs" // import generated docs
	bookingModels "myapp/modules/barberbooking/models"
	Core_controllers "myapp/modules/core/controllers"
	"myapp/modules/core/models"
	"myapp/modules/core/services"
	"myapp/modules/core/routes/admin"
    "myapp/modules/core/routes"
	"myapp/seeds"
)

func main() {
    app := fiber.New()

    // Global middleware
    app.Use(logger.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins:     "http://localhost:5173",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowCredentials: true,
    }))
    app.Use(recover.New())
    app.Use(helmet.New())
    app.Use(compress.New()) //บีบอัด response เพื่อลดขนาด

    // Connect & migrate
    database.ConnectDB()
    database.DB.AutoMigrate(
        &coreModels.Role{},
        &coreModels.Tenant{},
        &coreModels.Branch{},
        &coreModels.User{},
        &coreModels.TenantUser{},

        &bookingModels.Service{},
        &bookingModels.WorkingHour{},
        &bookingModels.Unavailability{},
        &bookingModels.Barber{},
        &bookingModels.Appointment{},
    )

     // 1) Seed Roles
     if err := seeds.SeedRoles(database.DB); err != nil {
        log.Fatalf("failed to seed roles: %v", err)
    }
    // 2) Seed Tenants
    if err := seeds.SeedTenants(database.DB); err != nil {
        log.Fatalf("failed to seed tenants: %v", err)
    }
    // 3) Seed Branch
    if err := seeds.SeedBranches(database.DB); err != nil {
        log.Fatalf("seed branches failed: %v", err)
    }
    // 4) Seed User
    if err := seeds.SeedUsers(database.DB); err != nil {
        log.Fatalf("seed users failed: %v", err)
    }
    // 5) Seed TenantUsers
    if err := seeds.SeedTenantUsers(database.DB); err != nil {
        log.Fatalf("seed tenant_users failed: %v", err)
    }
    // 6) Seed Services
    if err := seeds.SeedServices(database.DB); err != nil {
        log.Fatalf("seed services failed: %v", err)
    }
    // 7) Seed WorkingHours
    if err := seeds.SeedWorkingHours(database.DB); err != nil {
        log.Fatalf("seed working hours failed: %v", err)
    }
     // 8) Seed Unavailabilities
    if err := seeds.SeedUnavailabilities(database.DB); err != nil {
        log.Fatalf("seed unavailabilities failed: %v", err)
    }
      // 9) Seed Barbers
    if err := seeds.SeedBarbers(database.DB); err != nil {
        log.Fatalf("seed barbers failed: %v", err)
    }
      // 10) Seed Appointments
    if err := seeds.SeedAppointments(database.DB); err != nil {
        log.Fatalf("seed appointments failed: %v", err)
    }

    // Initialize Services & Controllers
    logSvc := services.NewSystemLogService(database.DB)
    Core_controllers.InitSystemLogHandler(logSvc)

    authSvc := services.NewAuthService(database.DB, logSvc)
    Core_controllers.InitAuthHandler(authSvc, logSvc)

    // Routes
    routes.SetupAuthRoutes(app)
    admin.SetupAdminRoutes(app)

    // Route api docs
    app.Get("/swagger/*", fiberSwagger.WrapHandler)

    // ลง middleware rate limiter หลัง route เพื่อจำกัดความถี่
    app.Use(limiter.New(limiter.Config{
        Max:        100,
        Expiration: 30 * time.Second,
    }))
    // ลอง deploy front-end 
    // app.Use("/", filesystem.New(filesystem.Config{
    //     Root:   http.Dir("/Users/nipatchapakdee/Mix_POS/frontend/dist"),
    //     Browse: false,
    //     Index:  "index.html",
    // }))

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "3001"
    }
    log.Fatal(app.Listen(":" + port))
}
