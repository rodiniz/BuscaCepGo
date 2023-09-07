package routes

import (
	"log"
	"strconv"
	"strings"

	"example.com/m/database"
	"example.com/m/models"
	"github.com/gofiber/fiber/v2"
)

// AllLocations is a function to get all locations from database
// @Summary Get all locations containing the name typed
// @Description Get all locations
// @Accept x-www-form-urlencoded
// @Param  name  formData  string  true  "Name"
// @Produce html
// @Success 200
// @Failure 500
// @Router /location [post]
func AllLocations(c *fiber.Ctx) error {
	locations := []models.Location{}
	nameFilter := c.FormValue("name")
	if err := database.DB.Db.Where("Logradouro LIKE ?", "%"+strings.ToUpper(nameFilter)+"%").Limit(15).Find(&locations).Error; err != nil {
		return c.Status(fiber.StatusOK).SendString("<tr><td colspan='3'>Rua não encontrada </td></tr>")
	}
	log.Println(len(locations))
	if len(locations) == 0 {
		return c.Status(fiber.StatusOK).SendString("<tr><td colspan='3'>Rua não encontrada </td></tr>")
	}
	builder := strings.Builder{}

	for i := 0; i < len(locations); i++ {
		builder.WriteString("<tr>")
		builder.WriteString("<td>" + locations[i].Uf + "</td>")
		builder.WriteString("<td>" + locations[i].Logradouro + "</td>")
		builder.WriteString("<td>" + locations[i].Cep + "</td>")
		builder.WriteString("</tr>")
	}
	return c.Status(200).SendString(builder.String())
}

// AllLocations is a function to get all locations from database
// @Summary Get all locations containing the name typed
// @Description Get all locations
// @Accept x-www-form-urlencoded
// @Param  name  query  string  true  "Name"
// @Param  page  query  string  true  "page"
// @Param  page_size  query  string  true  "page_size"
// @Produce json
// @Success 200
// @Failure 500
// @Router /location/paged [get]
func LocationsPaged(c *fiber.Ctx) error {
	var count int64
	locations := []models.Location{}
	nameFilter := c.Query("name")
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if err := database.DB.Db.Model(&locations).Scopes(database.LocationByName(nameFilter)).Count(&count).Error; err != nil {
		// Handle the error, print it, or log it
		log.Println("Error:", err)
	} else {
		log.Println("Count:", count)
	}
	if err := database.DB.Db.Scopes(database.LocationByName(nameFilter)).Scopes(database.Paginate(c)).Find(&locations).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(nil)
	}

	if len(locations) == 0 {
		return c.Status(fiber.StatusOK).JSON(nil)
	}
	totalPage := int(count/int64(pageSize)) + 1
	return c.Status(200).JSON(fiber.Map{
		"total":     count,
		"locations": locations,
		"page":      c.Query("page"),
		"nopages":   totalPage,
	})
}
