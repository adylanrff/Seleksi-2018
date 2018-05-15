package foodcrawler

import (
	"fmt"
	"regexp"
	"strconv"
)

/* Nutrient is a struct that represents the nutriuent contained in a food
 */

type Nutrient struct {
	Name  string  `json:"nutrient_name"`
	Unit  string  `json:"unit"`
	Value float32 `json:"value"`
}

func (nutrient Nutrient) String() string {
	return fmt.Sprintf("Name: %v\nUnit: %v\nValue: %v\n", nutrient.Name, nutrient.Unit, nutrient.Value)
}

/*Food is a struct that represents food in the database
 */
type Food struct {
	ID        string     `json:"id"`
	Name      string     `json:"food_name"`
	UPC       string     `json:"upc"`
	Nutrients []Nutrient `json:"nutrients"`
}

func (food Food) String() string {
	return fmt.Sprintf("ID:%v \nName:%v \nUPC: %v", food.ID, food.Name, food.UPC)
}

func (food Food) getDetailLink() string {
	detailLink := BaseURL + FoodDetailEndpoint + food.ID
	return detailLink
}

func parseFoodLink(foodText string) string {
	var productID, link string
	productIDPattern := regexp.MustCompile(`\d+(\s|,)`)
	productIDIdx := productIDPattern.FindStringSubmatchIndex(foodText)
	productID = foodText[productIDIdx[0]:(productIDIdx[1] - 1)]

	link = BaseURL + FoodDetailEndpoint + productID
	return link
}

func parseFoodData(foodText string, nutrients []Nutrient) Food {
	var upc, foodName, productID string
	productIDPattern := regexp.MustCompile(`\d+(\s|,)`)
	UPCPattern := regexp.MustCompile(`UPC:\s\d+`)

	productIDIdx := productIDPattern.FindStringSubmatchIndex(foodText)
	UPCIdx := UPCPattern.FindStringSubmatchIndex(foodText)

	productID = foodText[productIDIdx[0]:(productIDIdx[1] - 1)]

	if UPCIdx != nil {
		upc = foodText[(UPCIdx[0] + 5):UPCIdx[1]]
		foodName = foodText[productIDIdx[1]:(UPCIdx[0] - 2)]
	} else {
		foodName = foodText[productIDIdx[1]:]
	}

	foodData := Food{ID: productID, Name: foodName, UPC: upc, Nutrients: nutrients}

	return foodData
}

func parseNutrientData(name string, unit string, value string) Nutrient {
	var nutrient Nutrient
	stringPattern := regexp.MustCompile(`(\w+(\s)?[^(\n\r\t)]+)+`)

	name = stringPattern.FindString(name)

	floatValue, err := strconv.ParseFloat(value, 32)

	if err == nil {
		nutrient = Nutrient{Name: name, Unit: unit, Value: float32(floatValue)}
	}

	return nutrient
}
