package repositories

import (
	"github.com/PuerkitoBio/goquery"
)

type FeedbackRepository struct{}

func (s *FeedbackRepository) GetFeedbackPage(url string) (*goquery.Document, error) {
	doc, err := goquery.NewDocument(url)
	return doc, err
}
