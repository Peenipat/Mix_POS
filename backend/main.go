// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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

    _ "myapp/docs"              // import generated docs
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
    app.Use(compress.New()) //บีบอัด response เพื่อลดขนาด

    // Connect & migrate
    database.ConnectDB()
    database.DB.AutoMigrate(&models.User{}, &models.SystemLog{})

    // Initialize Services & Controllers
    logSvc := services.NewSystemLogService(database.DB)
    controllers.InitSystemLogHandler(logSvc)

    authSvc := services.NewAuthService(database.DB, logSvc)
    controllers.InitAuthHandler(authSvc, logSvc)

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
    app.Use("/", filesystem.New(filesystem.Config{
        Root:   http.Dir("/Users/nipatchapakdee/Mix_POS/frontend/dist"),
        Browse: false,
        Index:  "index.html",
    }))

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "3001"
    }
    log.Fatal(app.Listen(":" + port))
}
