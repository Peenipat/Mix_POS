package main
import (
    "github.com/gofiber/fiber/v2"
    "myapp/routes"
	"myapp/routes/admin"
    "myapp/database"
	"os"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"myapp/models"
	_ "myapp/docs"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)
func main(){
	// @title         Docs  api

// @host      localhost:3001
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,	}))
	

    database.ConnectDB()
	database.DB.AutoMigrate(&models.User{},&models.AuditLog{})
	routes.SetupAuthRoutes(app)
	admin.SetupAdminRoutes(app)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	if os.Getenv("ENV") != "production" {
		app.Get("/swagger/*", fiberSwagger.WrapHandler)
	}
	

	app.Listen(":3001")
}