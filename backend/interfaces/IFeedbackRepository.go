package interfaces

import (
	"github.com/PuerkitoBio/goquery"
)

type IFeedbackRepository interface {
	GetFeedbackPage(string) (*goquery.Document, error)
}
