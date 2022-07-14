package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rmnvlv/Web-Arctic/database"
	"github.com/rmnvlv/Web-Arctic/dynamic"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDatabase() {
	var err error
	database.DBconn, err = gorm.Open(sqlite.Open("test.db"))

	if err != nil {
		panic("Failed to connect to db")
	}
	fmt.Println("Succsess to connect to db")

	database.DBconn.AutoMigrate(&dynamic.Participant{})
	fmt.Println("DB migratetd")
}

// Init routes to get particions
func initDynamicRoutes(app *fiber.App) {

	admin := app.Group("/admin")
	dynamic.InitAdminRouter(admin)

	// admin.Use(jwtMiddleware.New(jwtMiddleware.Config{
	// 	ErrorHandler: func(c *fiber.Ctx, err error) error {
	// 		return fiber.NewError(1, "NOT_AUTHORIZED",
	// 			"Authentication credentials are missing or invalid.",
	// 			"Provide a properly configured and signed bearer token, and make sure that it has not expired.")
	// 	},
	// 	SigningMethod: jwt.SigningMethodHS256.Name,
	// 	SigningKey:    []byte(dynamic.JWT().SigningKey),
	// 	ContextKey:    "token",
	// }))

	user := app.Group("/user")
	dynamic.InitUserRoutes(user)

	// auth := app.Group("/auth")
	// dynamic.InitAuthRoutes(auth)
}

func main() {
	app := fiber.New()

	initDatabase()
	initDynamicRoutes(app)

	app.Listen(":8080")
}
