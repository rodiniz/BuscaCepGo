package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"example.com/m/database"
	_ "example.com/m/docs"
	"example.com/m/models"
	"example.com/m/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"
)

func main() {

	database.ConnectDb()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/location")
	api.Post("/", routes.AllLocations)
	api.Get("/paged", routes.LocationsPaged)
	app.Static("/img", "./img")
	app.Static("/public", "./public")
	app.Get("/", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("index", fiber.Map{}, "layouts/main")
	})

	//LoadData()

	app.Use(cors.New())
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	app.Listen(":" + port)
}

func LoadData() error {
	data, err := os.ReadFile("ceps.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}
	var locations []models.Location

	err2 := json.Unmarshal(data, &locations)
	if err2 != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err2
	}
	batchSize := 100 // You can adjust the batch size as needed

	for i := 0; i < len(locations); i += batchSize {
		end := i + batchSize
		if end > len(locations) {
			end = len(locations)
		}

		// Batch insert data into the database
		if err := database.DB.Db.Create(locations[i:end]).Error; err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("done running batches")
	return nil
}
