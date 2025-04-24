package main
import (
    "github.com/gofiber/fiber/v2"
    "myapp/routes"
    "myapp/database"
	"github.com/joho/godotenv"
	"os"
	"github.com/gofiber/fiber/v2/middleware/cors"
)
func main(){
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))
	
    database.ConnectDB()
	routes.SetupAuthRoutes(app)
	

	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load() 
	}

	app.Listen(":3001")
}