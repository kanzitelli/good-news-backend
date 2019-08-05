package crawler

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/kanzitelli/good-news-backend/models"
)

// SecretMag <struct>
// is used to present Secret Magazine crawler.
type SecretMag struct{}

const (
	baseURLSM  = "https://secretmag.ru"
	crawlURLSM = "https://secretmag.ru/news"
)

// Run <function>
// is used to start crawling process.
func (sm SecretMag) Run() []models.News {
	fmt.Println("Crawling SecretMag...")

	var totalNews []models.News
	newsFuncs := []NewsFunc{
		sm.runNews,
	}

	for _, f := range newsFuncs {
		tmpNews := f()
		totalNews = append(totalNews, tmpNews...)
	}

	return totalNews
}

func (sm SecretMag) runNews() []models.News {
	// creating simple colly instance without any options
	newsCollector := colly.NewCollector()

	// array of news that will be returned
	var news []models.News

	newsCollector.OnHTML(".wrapper", func(divWrapper *colly.HTMLElement) {
		divWrapper.ForEach(".container", func(i1 int, divContainer *colly.HTMLElement) {
			divContainer.ForEach(".item", func(i2 int, divItem *colly.HTMLElement) {
				link := divItem.ChildAttr("a[href]", "href")
				title := divItem.ChildText(".headline")

				news = append(news, models.News{
					Title:      title,
					Link:       fmt.Sprintf("%s%s", baseURLSM, link),
					Preamble:   "",
					TimeAdded:  time.Now().Unix(),
					NewsType:   models.TypeNews,
					NewsSource: models.SecretMagNewSource,
				})
			})
		})
	})

	newsCollector.Visit(crawlURLSM)
	newsCollector.Wait()

	return news
}
