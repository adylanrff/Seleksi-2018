package foodcrawler

import (
	"fmt"
	"regexp"
)

type Food struct {
	ID   string
	Name string
	UPC  string
}

type FoodDetail struct {
}

type Nutrient struct {
	Name string
	Unit string
}

func (food Food) String() string {
	return fmt.Sprintf("ID: %v\nName: %v\nUPC: %v", food.ID, food.Name, food.UPC)
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

func parseFoodData(foodText string) Food {
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

	foodData := Food{ID: productID, Name: foodName, UPC: upc}

	return foodData
}
