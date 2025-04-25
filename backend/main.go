package main
import (
    "github.com/gofiber/fiber/v2"
    "myapp/routes"
	"myapp/routes/admin"
    "myapp/database"
	"github.com/joho/godotenv"
	"os"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"myapp/models"
	"github.com/gofiber/swagger" 
	_ "myapp/docs"

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
		AllowCredentials: true,
	}))
	

    database.ConnectDB()
	database.DB.AutoMigrate(&models.User{})
	routes.SetupAuthRoutes(app)
	admin.SetupAdminRoutes(app)

	if os.Getenv("ENV") != "production" {
	app.Get("/swagger/*", swagger.HandlerDefault) // ใช้ใน dev mode เท่านั้น
}
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load() 
	}
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Listen(":3001")
}