package foodCrawler

import (
	"fmt"
	"os"
	"regexp"

	"github.com/anaskhan96/soup"
)

const BaseURL string = "https://ndb.nal.usda.gov/ndb"
const FoodsListEndpoint string = "/search/list"
const FoodDetailsEndpoint string = "/foods/show/"

func parseFoodList(foodText string) Food {
	productIDPattern := regexp.MustCompile(`\d+\s`)
	UPCPattern := regexp.MustCompile(`UPC:\s\d+`)

	productIDIdx := productIDPattern.FindStringSubmatchIndex(foodText)
	UPCIdx := UPCPattern.FindStringSubmatchIndex(foodText)

	productID := foodText[productIDIdx[0]:(productIDIdx[1] - 1)]
	upc := foodText[UPCIdx[0]:UPCIdx[1]]
	foodName := foodText[productIDIdx[1]:UPCIdx[0]]

	foodData := Food{ID: productID, Name: foodName, UPC: upc}

	return foodData
}

func CrawlBigPicture(linksList chan string, foodList chan Food, foodTexts []soup.Root) {

	for _, food := range foodTexts {
		foodLink := food.Find("a")
		foodData := parseFoodList(foodLink.Text())
		smallPictureLink := BaseURL + FoodDetailsEndpoint + foodData.ID

		linksList <- smallPictureLink
		foodList <- foodData
	}

	close(linksList)
	close(foodList)
}

func Crawl() {
	linksList := make(chan string)
	foodList := make(chan Food)
	resp, err := soup.Get(BaseURL + FoodsListEndpoint)
	if err != nil {
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)
	foodTexts := doc.FindAll("td", "style", "font-style:;")
	go CrawlBigPicture(linksList, foodList, foodTexts)

	for {
		link, ok1 := <-linksList
		food, ok2 := <-foodList

		if !ok1 || !ok2 {
			break
		} else {
			fmt.Println(link)
			fmt.Println(food)
		}

	}
	// fmt.Println(linksList[1].Text)
	// for _, food := range foodLists {
	// 	fmt.Println(food)
	// }
	// links := doc.Find("div", "id", "comicLinks").FindAll("a")
	// for _, link := range links {
	// 	linkElmt := Link{Text: link.Text(), Link: link.Attrs()["href"]}
	// 	linksList = append(linksList, linkElmt)
	// // }

	// out, _ := json.Marshal(linksList)
	// var linkz []Link

	// json.Unmarshal(out, &linkz)
	// fmt.Println(linkz[0].Link)
}
