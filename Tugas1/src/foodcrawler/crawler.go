package foodcrawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

/*BaseURL the base URL of the USDA Food Composition Databases
 */
const BaseURL string = "https://ndb.nal.usda.gov/ndb"

/*foodListEndpoint the endpoint of the list of foods
 */
const foodListEndpoint string = "/search/list?"

/*FoodDetailEndpoint the food detail endpoint that should be followed by the id of the food
 */
const FoodDetailEndpoint string = "/foods/show/"

func generateQuery(q string, offset int) string {
	return fmt.Sprintf("qlookup=%v&offset=%v", q, offset)
}

/*CrawlBigPicture crawls the big picture of the page.
Concurrently gets the small picture data using goroutine
*/
func CrawlBigPicture(linksChannel chan string, q string, count int) {
	for i := 0; i < int(count/25); i++ {
		query := generateQuery(q, i*25)
		resp, err := soup.Get(BaseURL + foodListEndpoint + query)
		if err != nil {
			os.Exit(1)
		}
		doc := soup.HTMLParse(resp)
		foodTexts := doc.FindAll("td", "style", "font-style:;")

		for _, food := range foodTexts {
			foodLink := food.Find("a")
			linksChannel <- parseFoodLink(foodLink.Text())
		}
	}

	close(linksChannel)
}

func parseTitle(doc soup.Root) (string, error) {
	title := doc.Find("div", "id", "view-name")
	if title.Error == nil {
		return title.Text(), nil
	}
	return "", errors.New("Can't find title")
}

func findValueColumn(header []soup.Root) int {
	for i := 0; i < len(header); i++ {
		if strings.Contains(header[i].Text(), "Value per 100 ") {
			return i
		}
	}
	return -1
}

func parseNutrients(nutlist soup.Root, index int) []Nutrient {

	var nutrients []Nutrient
	rows := nutlist.FindAll("tr")
	for _, val := range rows {

		var newNutrient Nutrient
		if val.Attrs()["class"] != "group" {
			columns := val.FindAll("td")
			if len(columns) > 0 {
				nutrient := columns[1].Text()
				unit := columns[2].Text()
				value := columns[index].Text()

				newNutrient = parseNutrientData(nutrient, unit, value)
				nutrients = append(nutrients, newNutrient)

			}
		}
	}
	return nutrients
}

func parseNutritionTable(doc soup.Root) ([]Nutrient, error) {
	var nutrients []Nutrient
	nutritionList := doc.Find("div", "class", "nutlist")
	if nutritionList.Error == nil {

		// Parse header and get the value of the column
		header := nutritionList.FindAll("th")
		valueIdx := findValueColumn(header)
		if valueIdx == -1 {
			return nutrients, errors.New("Can't find value table")
		}
		nutrients = parseNutrients(nutritionList, valueIdx)

		return nutrients, nil
	}

	return nutrients, errors.New("Can't find nutrition table")
}

/*CrawlSmallPicture crawls the details of the food
Returns the foodDetail struct
*/
func CrawlSmallPicture(linksChannel chan string, foodChannel chan Food) {
	for {
		link, ok := <-linksChannel
		if !ok {
			break
		}

		resp, _ := soup.Get(link)
		doc := soup.HTMLParse(resp)
		if doc.Error == nil {
			foodTitle, errTitle := parseTitle(doc)
			nutritionTable, errTable := parseNutritionTable(doc)

			if errTitle == nil && errTable == nil {
				foodData := parseFoodData(foodTitle, nutritionTable)
				foodChannel <- foodData
			}

		}

		time.Sleep(time.Millisecond * 1000)
	}

	close(foodChannel)
}

/*Crawl crawl the USDA website
 */
func Crawl(filename string, q string, count int) []Food {

	linksChannel := make(chan string)
	foodChannel := make(chan Food)

	var foodList []Food

	go CrawlBigPicture(linksChannel, q, count)
	go CrawlSmallPicture(linksChannel, foodChannel)

	for i := 1; ; i++ {
		food, ok := <-foodChannel
		if !ok {
			break
		} else {
			foodList = append(foodList, food)
			fmt.Println("Link-", i, " parsed")
		}
	}

	out, _ := json.Marshal(foodList)
	err := ioutil.WriteFile(filename, out, 0644)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return foodList
}
