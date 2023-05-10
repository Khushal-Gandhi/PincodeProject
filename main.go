package main

import (
	"encoding/csv"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"os"
)

type Holder struct {
	Map map[string][]map[string]string
}

func main() {
	path := "Locality_village_pincode_final_mar-2017.csv"

	h := Holder{}
	Mapper := make(map[string][]map[string]string)
	h.Map = Mapper

	h.loadCSV(path)

	app := fiber.New()

	app.Get("/api/search/:pincode", func(c *fiber.Ctx) error {

		pincode := c.Params("pincode")

		data, err := h.GetPincodeData(pincode)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not get pincode data",
			})
		}

		return c.JSON(data)
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start the server %v", err)
	}
}

func (h *Holder) loadCSV(path string) {

	//    read csv
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	//defer f.Close()
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("Failed to close file: %v", err)
		}
	}()
	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, datum := range data {

		pincode := datum[2]
		state := datum[5]
		area := datum[1]
		city := datum[4]

		InnerMapper := make(map[string]string)

		InnerMapper["pincode"] = pincode
		InnerMapper["state"] = state
		InnerMapper["area"] = area
		InnerMapper["city"] = city

		h.Map[pincode] = append(h.Map[pincode], InnerMapper)

	}

}

func (h *Holder) GetPincodeData(pincode string) ([]map[string]string, error) {

	return h.Map[pincode], nil
}
