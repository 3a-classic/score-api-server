package route

import (
	"logger"
	"mongo"

	"log"
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

	logger.Output(
		logrus.Fields{
			"Page":     page,
			"Team":     team,
			"Hole":     hole,
			"Request":  req,
			"Response": resp,
		},
		"Post access to page router",
		logger.Debug,
	)
	switch page {

	case "login":

		loginInfo := new(mongo.PostLogin)
		err := req.ReadEntity(loginInfo)
		if err != nil { // bad request
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := mongo.PostLoginPageData(loginInfo)
		if err != nil {
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:   err,
					logger.TraceMsg: logger.Trace(),
					"Status":        status,
					"login Info":    loginInfo,
				},
				"PostLoginPageData Err",
				logger.Error,
			)
		}
		resp.WriteAsJson(status)

	case "scoreViewSheet":

		definedTeam := new(mongo.PostDefinedTeam)
		err := req.ReadEntity(definedTeam)
		if err != nil { // bad request
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := mongo.PostScoreViewSheetPageData(team, definedTeam)
		if err != nil {
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:   err,
					logger.TraceMsg: logger.Trace(),
					"Status":        status,
					"Team":          team,
					"Defind Team":   definedTeam,
				},
				"PostScoreViewSheetPageData Err",
				logger.Error,
			)
		}
		resp.WriteAsJson(status)

	case "scoreEntrySheet":

		updatedTeamScore := new(mongo.PostTeamScore)
		err := req.ReadEntity(updatedTeamScore)
		log.Println(updatedTeamScore)
		if err != nil { // bad request
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := mongo.PostScoreEntrySheetPageData(team, hole, updatedTeamScore)
		if err != nil {
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:        err,
					logger.TraceMsg:      logger.Trace(),
					"Status":             status,
					"Team":               team,
					"Hole":               hole,
					"Updated Team Score": updatedTeamScore,
				},
				"PostScoreEntrySheetPageData Err",
				logger.Error,
			)
		}
		resp.WriteAsJson(status)

	case "applyScore":
		if hole != "" {
			return
		}
		registeredApplyScore := new(mongo.PostApplyScore)
		err := req.ReadEntity(registeredApplyScore)
		if err != nil {
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := mongo.PostApplyScoreData(team, registeredApplyScore)
		if err != nil {
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:   err,
					logger.TraceMsg: logger.Trace(),
					"Status":        status,
					"Team":          team,
					"Registered Apply Score": registeredApplyScore,
				},
				"PostApplyScoreData Err",
				logger.Error,
			)
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

	logger.Output(
		logrus.Fields{
			"Date":       date,
			"Collection": collection,
			"Request":    req,
			"Response":   resp,
		},
		"Post access to register router",
		logger.Debug,
	)

	switch collection {

	case "user":

		userCols := new([]mongo.UserCol)
		log.Println(userCols)
		if err := req.ReadEntity(userCols); err != nil {
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := mongo.RegisterUserColData(*userCols)
		if err != nil {
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:      err,
					logger.TraceMsg:    logger.Trace(),
					"Status":           status,
					"User Collections": userCols,
				},
				"RegisterUserColData Err",
				logger.Error,
			)
		}
		resp.WriteAsJson(status)

	case "team":

		teamCols := new([]mongo.TeamCol)
		log.Println(teamCols)
		if err := req.ReadEntity(teamCols); err != nil {
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := mongo.RegisterTeamColData(date, *teamCols)
		if err != nil {
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:      err,
					logger.TraceMsg:    logger.Trace(),
					"Status":           status,
					"Date":             date,
					"Team Collections": teamCols,
				},
				"RegisterTeamColData Err",
				logger.Error,
			)
		}
		resp.WriteAsJson(status)

	case "field":

		fieldCols := new([]mongo.FieldCol)
		log.Println(fieldCols)
		if err := req.ReadEntity(fieldCols); err != nil {
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := mongo.RegisterFieldColData(date, *fieldCols)
		if err != nil {
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:       err,
					logger.TraceMsg:     logger.Trace(),
					"Status":            status,
					"Date":              date,
					"Field Collections": fieldCols,
				},
				"RegisterFieldColData Err",
				logger.Error,
			)
		}
		resp.WriteAsJson(status)

	case "thread":

		requestTakePictureStatus := new(mongo.RequestTakePictureStatus)
		log.Println(requestTakePictureStatus)
		if err := req.ReadEntity(requestTakePictureStatus); err != nil {
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := mongo.RegisterThreadImg(requestTakePictureStatus)
		if err != nil {
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:                 err,
					logger.TraceMsg:               logger.Trace(),
					"Status":                      status,
					"Request Take Picture Status": requestTakePictureStatus,
				},
				"RegisterThreadImg Err",
				logger.Error,
			)
		}

		resp.WriteAsJson(status)
	}
}
