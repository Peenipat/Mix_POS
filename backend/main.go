package main

import (
    "os"

    "myapp/database"
    "myapp/models"
    "myapp/routes"
    "myapp/routes/admin"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/gofiber/fiber/v2/middleware/helmet"
    "github.com/gofiber/fiber/v2/middleware/compress"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "time"
    //deploy
    "github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"

    _ "myapp/docs" // Swagger docs
    fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title         Docs API
// @host          localhost:3001
// @BasePath      /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
    app := fiber.New()

    // Middleware
    app.Use(logger.New())
    app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", 
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))
	
    app.Use(recover.New())
    app.Use(helmet.New())
    app.Use(compress.New())

    // Connect Database
    database.ConnectDB()
    database.DB.AutoMigrate(&models.User{}, &models.AuditLog{})

    // Swagger Docs
    app.Get("/swagger/*", fiberSwagger.WrapHandler)

    // API Routes
    routes.SetupAuthRoutes(app)
    admin.SetupAdminRoutes(app)

    app.Use(limiter.New(limiter.Config{
        Max:        100,                  // เพิ่มจำนวน request ที่อนุญาต
        Expiration: 30 * time.Second,     // ขยายเวลา reset
    }))

    // Serve React Static Files (หลังสุด)
    app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.Dir("/Users/nipatchapakdee/Mix_POS/frontend/dist"), // ใช้ net/http เลย
		Browse:     false,
		Index:      "index.html",
	}))

    // Listen
    port := os.Getenv("PORT")
    if port == "" {
        port = "3001"
    }
    app.Listen(":" + port)
}
