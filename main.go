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
	data, err := os.ReadFile("C:\\rodrigo\\ceps.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}
	var myData []models.Location

	err2 := json.Unmarshal(data, &myData)
	if err2 != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err2
	}
	startIndex := 100000
	fmt.Println(len(myData))
	fmt.Println("Populating array")
	var locations = []*models.Location{}
	for i := startIndex; i < 110000; i++ {
		v := myData[i]

		if !exists(v.Cep) {
			fmt.Println(i)
			book := &models.Location{
				Bairro:         v.Bairro,
				Cep:            v.Cep,
				TipoLogradouro: v.TipoLogradouro,
				Logradouro:     v.Logradouro,
				Cidade:         v.Cidade,
				Uf:             v.Uf,
				CodigoIbge:     v.CodigoIbge,
			}
			locations = append(locations, book)
		}
	}
	fmt.Println("Running batches")
	result := database.DB.Db.CreateInBatches(locations, 400)
	if result.Error != nil {
		fmt.Println("Error", result.Error)
		return err2
	}
	fmt.Println("done running batches")
	return nil
}
func exists(cep string) bool {
	var location models.Location
	result := database.DB.Db.Where(&models.Location{Cep: cep}).First(&location)
	return result.Error == nil
}
