// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"os"
	"time"
    "github.com/joho/godotenv"

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
	
	Core_controllers "myapp/modules/core/controllers"	
	"myapp/modules/core/services"
	"myapp/modules/core/routes/admin"
    "myapp/modules/core/routes"
	
)

func main() {
    // Connect & migrate
    database.ConnectDB()

    app := fiber.New()

    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, relying on real environment variables")
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("Missing JWT_SECRET")
    }

    app.Use(logger.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins:     "https://nipat-cv-com-cp.vercel.app, http://localhost:5173",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowCredentials: true,
    }))
    app.Use(recover.New())
    app.Use(helmet.New())
    app.Use(compress.New()) 
    // Initialize Services & Controllers
    
    logSvc := services.NewSystemLogService(database.DB)
    Core_controllers.InitSystemLogHandler(logSvc)

    authSvc := services.NewAuthService(database.DB, logSvc)
    Core_controllers.InitAuthHandler(authSvc, logSvc)

    // Routes
    routes.SetupAuthRoutes(app)
    admin.SetupAdminRoutes(app)

    // Global middleware
    //บีบอัด response เพื่อลดขนาด



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
