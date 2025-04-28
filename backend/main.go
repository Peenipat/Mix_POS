// main.go
package main

import (
    "log"
    "os"
    "time"
    "net/http"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/gofiber/fiber/v2/middleware/helmet"
    "github.com/gofiber/fiber/v2/middleware/compress"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/fiber/v2/middleware/filesystem"
    fiberSwagger "github.com/swaggo/fiber-swagger"

    "myapp/controllers"
    "myapp/database"
    "myapp/models"
    "myapp/routes"
    "myapp/routes/admin"
    "myapp/services"
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
    app.Use(compress.New())

    // Connect & migrate
    database.ConnectDB()
    database.DB.AutoMigrate(&models.User{}, &models.SystemLog{})

    // Services & Controllers initialization
    logSvc := services.NewSystemLogService(database.DB)
    controllers.InitSystemLogHandler(logSvc)

    authSvc := services.NewAuthService(database.DB, logSvc)
    controllers.InitAuthHandler(authSvc, logSvc)

    // Swagger
    app.Get("/swagger/*", fiberSwagger.WrapHandler)

    // Routes
    routes.SetupAuthRoutes(app)
    admin.SetupAdminRoutes(app)

    // Rate limiter & static
    app.Use(limiter.New(limiter.Config{
        Max:        100,
        Expiration: 30 * time.Second,
    }))
    app.Use("/", filesystem.New(filesystem.Config{
        Root:   http.Dir("/Users/nipatchapakdee/Mix_POS/frontend/dist"),
        Browse: false,
        Index:  "index.html",
    }))

    // Start
    port := os.Getenv("PORT")
    if port == "" {
        port = "3001"
    }
    log.Fatal(app.Listen(":" + port))
}
