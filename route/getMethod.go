package route

import (
	l "github.com/3a-classic/score-api-server/logger"
	m "github.com/3a-classic/score-api-server/mongo"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
)

func getCol(c *gin.Context) {
	col := c.Params.ByName("col")

	l.Output(
		logrus.Fields{"Collection": l.Sprintf(col)},
		"Get access to collection router",
		l.Debug,
	)

	if col == "player" || col == "field" || col == "team" {
		data, err := m.GetAllColData(col)
		pp.Println(data)
		if err != nil {
			panic(err)
		}
		c.JSON(200, data)
	}
}

func getPage(c *gin.Context) {
	var page, team, hole string
	url := strings.Split(c.Params.ByName("page"), "/")
	pp.Println(url)

	l.Output(
		logrus.Fields{"Page": page, "Team": team, "Hole": hole, "URL": url},
		"Get access to page router",
		l.Debug,
	)
	page = url[1]
	if len(url) > 2 {
		team = url[2]
	}
	if len(url) > 3 {
		hole = url[3]
	}

	pp.Println(page)
	switch page {
	case "index":
		if team != "" || hole != "" {
			return
		}
		data, err := m.GetIndexPageData()
		if err != nil {
			panic(err)
		}
		c.JSON(200, data)

	case "leadersBoard":
		if team != "" || hole != "" {
			return
		}
		data, err := m.GetLeadersBoardPageData()
		if err != nil {
			panic(err)
		}
		c.JSON(200, data)

	case "scoreEntrySheet":
		data, err := m.GetScoreEntrySheetPageData(team, hole)
		if err != nil {
			panic(err)
		}
		c.JSON(200, data)

	case "scoreViewSheet":
		if hole != "" {
			return
		}
		data, err := m.GetScoreViewSheetPageData(team)
		if err != nil {
			panic(err)
		}
		c.JSON(200, data)

	case "entireScore":

		data, err := m.GetEntireScorePageData()
		if err != nil {
			panic(err)
		}
		c.JSON(200, data)

	case "timeLine":

		data, err := m.GetTimeLinePageData()
		if err != nil {
			panic(err)
		}
		c.JSON(200, data)

	default:
		c.JSON(404, gin.H{"status": "not found"})
	}
}
