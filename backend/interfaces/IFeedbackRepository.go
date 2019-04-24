package interfaces

import (
	"github.com/PuerkitoBio/goquery"
)

type IFeedbackRepository interface {
	GetFeedbackPage(string) (*goquery.Document, int, error)
}
