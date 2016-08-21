package route

import (
	l "logger"
	m "mongo"

	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

func postOne(req *restful.Request, resp *restful.Response) {
	page := req.PathParameter("page")
	team := req.PathParameter("team")
	hole := req.PathParameter("hole")
	if origin := req.HeaderParameter(restful.HEADER_Origin); origin != "" {
		if len(resp.Header().Get(restful.HEADER_AccessControlAllowOrigin)) == 0 {
			resp.AddHeader(restful.HEADER_AccessControlAllowOrigin, origin)
		}
	}

	l.Output(
		logrus.Fields{
			"Page":     page,
			"Team":     team,
			"Hole":     hole,
			"Request":  l.Sprintf(req),
			"Response": l.Sprintf(resp),
		},
		"Post access to page router",
		l.Debug,
	)
	switch page {

	case "login":

		loginInfo := new(m.PostLogin)
		err := req.ReadEntity(loginInfo)
		if err != nil { // bad request
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := m.PostLoginPageData(loginInfo)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_PostPage, loginInfo)
		}
		resp.WriteAsJson(status)

	case "scoreViewSheet":

		definedTeam := new(m.PostDefinedTeam)
		err := req.ReadEntity(definedTeam)
		if err != nil { // bad request
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := m.PostScoreViewSheetPageData(team, definedTeam)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_PostPage, definedTeam)
		}
		resp.WriteAsJson(status)

	case "scoreEntrySheet":

		updatedTeamScore := new(m.PostTeamScore)
		err := req.ReadEntity(updatedTeamScore)
		if err != nil { // bad request
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := m.PostScoreEntrySheetPageData(team, hole, updatedTeamScore)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_PostPage, updatedTeamScore)
		}
		resp.WriteAsJson(status)

	case "applyScore":
		if hole != "" {
			return
		}
		registeredApplyScore := new(m.PostApplyScore)
		err := req.ReadEntity(registeredApplyScore)
		if err != nil {
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := m.PostApplyScoreData(team, registeredApplyScore)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_PostPage, registeredApplyScore)
		}
		resp.WriteAsJson(status)
	}
}

func register(req *restful.Request, resp *restful.Response) {
	date := req.PathParameter("date")
	collection := req.PathParameter("collection")
	if origin := req.HeaderParameter(restful.HEADER_Origin); origin != "" {
		if len(resp.Header().Get(restful.HEADER_AccessControlAllowOrigin)) == 0 {
			resp.AddHeader(restful.HEADER_AccessControlAllowOrigin, origin)
		}
	}

	l.Output(
		logrus.Fields{
			"Date":       date,
			"Collection": collection,
			"Request":    l.Sprintf(req),
			"Response":   l.Sprintf(resp),
		},
		"Post access to register router",
		l.Debug,
	)

	switch collection {

	case "user":

		userCols := new([]m.UserCol)
		if err := req.ReadEntity(userCols); err != nil {
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := m.RegisterUserColData(*userCols)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_RegisterCol, userCols)
		}
		resp.WriteAsJson(status)

	case "team":

		teamCols := new([]m.TeamCol)
		if err := req.ReadEntity(teamCols); err != nil {
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := m.RegisterTeamColData(date, *teamCols)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_RegisterCol, teamCols)
		}
		resp.WriteAsJson(status)

	case "field":

		fieldCols := new([]m.FieldCol)
		if err := req.ReadEntity(fieldCols); err != nil {
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := m.RegisterFieldColData(date, *fieldCols)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_RegisterCol, fieldCols)
		}
		resp.WriteAsJson(status)

	case "thread":

		requestTakePictureStatus := new(m.RequestTakePictureStatus)
		if err := req.ReadEntity(requestTakePictureStatus); err != nil {
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := m.RegisterThreadImg(requestTakePictureStatus)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_RegisterCol, requestTakePictureStatus)
		}

		resp.WriteAsJson(status)
	}
}
