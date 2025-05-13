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
	bookingControllers "myapp/modules/barberbooking/controllers"
	bookingModels "myapp/modules/barberbooking/models"
	bookingRoutes "myapp/modules/barberbooking/routes"
	bookingServices "myapp/modules/barberbooking/services"
	Core_controllers "myapp/modules/core/controllers"
	"myapp/modules/core/models"
	coreServices "myapp/modules/core/services"
	"myapp/routes"
	"myapp/routes/admin"
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
	if database.DB == nil {
		log.Fatal("❌ GORM DB is nil. Cannot proceed.")
	}
	
	database.DB.AutoMigrate(
		// Core module: สร้างสิ่งที่เป็นรากก่อน
		&coreModels.Tenant{},
		&coreModels.Role{},
		&coreModels.Module{},
		&coreModels.Branch{},
		&coreModels.User{},
		&coreModels.TenantUser{},
		&coreModels.TenantModule{},

		// Booking module
		&bookingModels.Customer{},
		&bookingModels.Service{},
		&bookingModels.WorkingHour{},
		&bookingModels.Barber{},
		&bookingModels.Unavailability{},
		&bookingModels.Appointment{},
		&bookingModels.AppointmentStatusLog{},
		&bookingModels.AppointmentReview{},
		&bookingModels.BarberWorkload{},
	)

	// 1) Seed Tenants → เพื่อให้มี tenant ใช้ใน Role, Branch, User
	if err := seeds.SeedTenants(database.DB); err != nil {
		log.Fatalf("failed to seed tenants: %v", err)
	}

	// 2) Seed Modules → ระบบรองรับ feature ของ tenant (tenant_modules จะตามมาภายหลัง)
	if err := seeds.SeedModules(database.DB); err != nil {
		log.Fatalf("seed modules failed: %v", err)
	}

	// 3) Seed Roles → ต้องใช้ TenantID (บาง role อาจเป็น per-tenant)
	if err := seeds.SeedRoles(database.DB); err != nil {
		log.Fatalf("failed to seed roles: %v", err)
	}

	// 4) Seed Branches → ใช้ tenant_id
	if err := seeds.SeedBranches(database.DB); err != nil {
		log.Fatalf("seed branches failed: %v", err)
	}

	// 5) Seed Users → ใช้ role_id และ branch_id
	if err := seeds.SeedUsers(database.DB); err != nil {
		log.Fatalf("seed users failed: %v", err)
	}

	// 6) Seed TenantUsers → ต้องมี tenant และ user
	if err := seeds.SeedTenantUsers(database.DB); err != nil {
		log.Fatalf("seed tenant_users failed: %v", err)
	}

	tenantID := uint(1)
	// 7) Seed Customers → เป็นลูกค้าจากภายนอก ไม่ต้องพึ่ง tenant_id
	if err := seeds.SeedCustomers(database.DB,tenantID); err != nil {
		log.Fatalf("seed customers failed: %v", err)
	}

	// 8) Seed Services → ข้อมูลภายใน barber module
	if err := seeds.SeedServices(database.DB); err != nil {
		log.Fatalf("seed services failed: %v", err)
	}

	// 9) Seed WorkingHours → ต้องมี branch
	if err := seeds.SeedWorkingHours(database.DB); err != nil {
		log.Fatalf("seed working hours failed: %v", err)
	}

	// 10) Seed Unavailabilities → ต้องมี branch และ (optional) barber
	if err := seeds.SeedUnavailabilities(database.DB); err != nil {
		log.Fatalf("seed unavailabilities failed: %v", err)
	}

	// 11) Seed Barbers → ต้องมี branch และ user
	if err := seeds.SeedBarbers(database.DB); err != nil {
		log.Fatalf("seed barbers failed: %v", err)
	}

	// 12) Seed Appointments → ต้องมี branch, service, barber, customer
	if err := seeds.SeedAppointments(database.DB); err != nil {
		log.Fatalf("seed appointments failed: %v", err)
	}

	// 13) Seed Appointment Status Logs → ต้องมี appointment
	if err := seeds.SeedAppointmentStatusLogs(database.DB); err != nil {
		log.Fatalf("seed appointment status logs failed: %v", err)
	}

	// 14) Seed Appointment Reviews → ต้องมี appointment และ customer
	if err := seeds.SeedAppointmentReviews(database.DB); err != nil {
		log.Fatalf("seed appointment reviews failed: %v", err)
	}

	// 15) Seed Barber Workloads → ต้องมี barbers + appointments
	if err := seeds.SeedBarberWorkloads(database.DB); err != nil {
		log.Fatalf("seed barber workloads failed: %v", err)
	}

	// 16) Seed TenantModules
	if err := seeds.SeedTenantModules(database.DB); err != nil {
		log.Fatalf("seed tenant modules failed: %v", err)
	}

	routes.SetupAuthRoutes(app)
	admin.SetupAdminRoutes(app)
	// Initialize Services & Controllers
	logSvc := coreServices.NewSystemLogService(database.DB)
	Core_controllers.InitSystemLogHandler(logSvc)

	authSvc := coreServices.NewAuthService(database.DB, logSvc)
	Core_controllers.InitAuthHandler(authSvc, logSvc)

	// === Barber Booking Module: Service Feature ===
	serviceService := bookingServices.NewServiceService(database.DB)
	serviceController := bookingControllers.NewServiceController(serviceService)

	customerService := bookingServices.NewCustomerService(database.DB)
	customerController := bookingControllers.NewCustomerController(customerService) 

	barberService := bookingServices.NewBarberService(database.DB)
	barberController := bookingControllers.NewBarberController(barberService)

	unavailabilityService := bookingServices.NewUnavailabilityService(database.DB)
	unavailabilityController := bookingControllers.NewUnavailabilityController(unavailabilityService)
	
	workingHourService := bookingServices.NewWorkingHourService(database.DB)
	workingHourController := bookingControllers.NewWorkingHourController(workingHourService)

	bookingGroup := app.Group("/api/barberbooking")

	// Register routes
	bookingRoutes.RegisterServiceRoutes(bookingGroup, serviceController)
	bookingRoutes.RegisterCustomerRoutes(bookingGroup,customerController)
	bookingRoutes.RegisterBarberRoutes(bookingGroup,barberController)
	bookingRoutes.RegisterUnavailabilityRoute(bookingGroup,unavailabilityController)
	bookingRoutes.RegisterWorkingHourRoute(bookingGroup,*workingHourController)
	

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
