package foodcrawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/anaskhan96/soup"
)

/*BaseURL the base URL of the USDA Food Composition Databases
 */
const BaseURL string = "https://ndb.nal.usda.gov/ndb"

/*foodListEndpoint the endpoint of the list of foods
 */
const foodListEndpoint string = "/search/list"

/*FoodDetailEndpoint the food detail endpoint that should be followed by the id of the food
 */
const FoodDetailEndpoint string = "/foods/show/"

/*CrawlBigPicture crawls the big picture of the page.
Concurrently gets the small picture data using goroutine
*/
func CrawlBigPicture(linksChannel chan string) {

	resp, err := soup.Get(BaseURL + foodListEndpoint)
	if err != nil {
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)
	foodTexts := doc.FindAll("td", "style", "font-style:;")

	for _, food := range foodTexts {
		foodLink := food.Find("a")
		linksChannel <- parseFoodLink(foodLink.Text())
	}

	close(linksChannel)
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
			foodTexts := doc.Find("div", "id", "view-name")
			if foodTexts.Error == nil {
				foodData := parseFoodData(foodTexts.Text())
				foodChannel <- foodData
			}
		}
	}

	close(foodChannel)
}

/*Crawl crawl the USDA website
 */
func Crawl() {

	linksChannel := make(chan string)
	foodChannel := make(chan Food)

	var foodList []Food

	go CrawlBigPicture(linksChannel)
	go CrawlSmallPicture(linksChannel, foodChannel)

	for i := 1; ; i++ {
		food, ok := <-foodChannel
		if !ok {
			break
		} else {
			fmt.Println("Link-", i, " parsed")
			foodList = append(foodList, food)
		}
	}

	out, _ := json.Marshal(foodList)
	err := ioutil.WriteFile("output.json", out, 0644)

	if err != nil {
		os.Exit(1)
	}
}
