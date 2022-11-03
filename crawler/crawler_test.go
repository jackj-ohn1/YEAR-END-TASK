package crawler

import (
	"context"
	"testing"
	"year-end/model"
)

func TestCrawler(t *testing.T) {
	model.D.Init()
	model.DR.Init()
	Crawler(context.Background(), "", "")
}
