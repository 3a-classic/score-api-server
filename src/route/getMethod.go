package route

import (
	"logger"
	"mongo"

	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

func getCol(req *restful.Request, resp *restful.Response) {
	col := req.PathParameter("col")

	logger.Output(
		logrus.Fields{
			"Collection": col,
		},
		"Get access to collection router",
		logger.Debug,
	)

	if col == "player" || col == "field" || col == "team" {
		data, err := mongo.GetAllColData(col)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
	}
}

func getPage(req *restful.Request, resp *restful.Response) {
	var page, team, hole string
	url := (strings.Split(req.PathParameter("page"), "/"))

	logger.Output(
		logrus.Fields{
			"Page": page,
			"Team": team,
			"Hole": hole,
			"URL":  url,
		},
		"Get access to page router",
		logger.Debug,
	)
	page = url[0]
	if len(url) > 1 {
		team = url[1]
	}
	if len(url) > 2 {
		hole = url[2]
	}

	switch page {
	case "index":
		if team != "" || hole != "" {
			return
		}
		data, err := mongo.GetIndexPageData()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)

	case "leadersBoard":
		if team != "" || hole != "" {
			return
		}
		data, err := mongo.GetLeadersBoardPageData()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)

	case "scoreEntrySheet":
		data, err := mongo.GetScoreEntrySheetPageData(team, hole)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)

	case "scoreViewSheet":
		if hole != "" {
			return
		}
		data, err := mongo.GetScoreViewSheetPageData(team)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)

	case "entireScore":

		data, err := mongo.GetEntireScorePageData()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)

	case "timeLine":

		data, err := mongo.GetTimeLinePageData()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)

	default:
		resp.WriteErrorString(
			http.StatusNotFound,
			"404: Page is not found.",
		)
	}
}
