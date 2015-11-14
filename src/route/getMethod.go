package route

import (
	l "logger"
	m "mongo"

	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

func getCol(req *restful.Request, resp *restful.Response) {
	col := req.PathParameter("col")

	l.Output(
		logrus.Fields{"Collection": col},
		"Get access to collection router",
		l.Debug,
	)

	if col == "player" || col == "field" || col == "team" {
		data, err := m.GetAllColData(col)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
	}
}

func getPage(req *restful.Request, resp *restful.Response) {
	var page, team, hole string
	url := (strings.Split(req.PathParameter("page"), "/"))

	l.Output(
		logrus.Fields{"Page": page, "Team": team, "Hole": hole, "URL": url},
		"Get access to page router",
		l.Debug,
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
		data, err := m.GetIndexPageData()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)

	case "leadersBoard":
		if team != "" || hole != "" {
			return
		}
		data, err := m.GetLeadersBoardPageData()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)

	case "scoreEntrySheet":
		data, err := m.GetScoreEntrySheetPageData(team, hole)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)

	case "scoreViewSheet":
		if hole != "" {
			return
		}
		data, err := m.GetScoreViewSheetPageData(team)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)

	case "entireScore":

		data, err := m.GetEntireScorePageData()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)

	case "timeLine":

		data, err := m.GetTimeLinePageData()
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
