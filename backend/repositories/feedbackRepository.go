package repositories

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

type FeedbackRepository struct{}

func (s *FeedbackRepository) GetFeedbackPage(url string) (*goquery.Document, int, error) {
	code			:= 666
	response, _ := http.Get(url)
	doc, err    := goquery.NewDocument(url)

	if response != nil {
		code = response.StatusCode
	}

	return doc, code, err
}