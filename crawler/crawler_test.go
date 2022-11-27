package crawler

import (
	"testing"
	"year-end/config"
	"year-end/model"
)

func TestCrawler(t *testing.T) {
	model.D.Init()
	model.DR.Init()
	cfg := &config.Config{}
	cfg.Init("./config/config.yaml")
}
