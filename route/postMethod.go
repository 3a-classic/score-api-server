package route

import (
	l "github.com/3a-classic/score-api-server/logger"
	m "github.com/3a-classic/score-api-server/mongo"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
)

func postOne(c *gin.Context) {
	page := c.Params.ByName("page")
	team := c.Params.ByName("team")
	hole := c.Params.ByName("hole")
	//		if origin := req.HeaderParameter(restful.HEADER_Origin); origin != "" {
	//			if len(resp.Header().Get(restful.HEADER_AccessControlAllowOrigin)) == 0 {
	//				resp.AddHeader(restful.HEADER_AccessControlAllowOrigin, origin)
	//			}
	//		}

	l.Output(
		logrus.Fields{
			"Page":    page,
			"Team":    team,
			"Hole":    hole,
			"Context": l.Sprintf(c),
		},
		"Post access to page router",
		l.Debug,
	)
	switch page {

	case "login":

		loginInfo := new(m.PostLogin)
		if !c.Bind(loginInfo) {
			c.JSON(404, gin.H{"status": "wrong queries"})
			return
		}

		status, err := m.PostLoginPageData(loginInfo)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_PostPage, loginInfo)
		}
		c.JSON(200, status)

	case "scoreViewSheet":

		definedTeam := new(m.PostDefinedTeam)
		if !c.Bind(definedTeam) {
			c.JSON(404, gin.H{"status": "wrong queries"})
			return
		}

		status, err := m.PostScoreViewSheetPageData(team, definedTeam)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_PostPage, definedTeam)
		}
		c.JSON(200, status)

	case "scoreEntrySheet":

		updatedTeamScore := new(m.PostTeamScore)
		if !c.Bind(updatedTeamScore) {
			c.JSON(404, gin.H{"status": "wrong queries"})
			return
		}

		status, err := m.PostScoreEntrySheetPageData(team, hole, updatedTeamScore)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_PostPage, updatedTeamScore)
		}
		c.JSON(200, status)

	case "applyScore":
		if hole != "" {
			return
		}
		registeredApplyScore := new(m.PostApplyScore)
		if !c.Bind(registeredApplyScore) {
			c.JSON(404, gin.H{"status": "wrong queries"})
			return
		}

		status, err := m.PostApplyScoreData(team, registeredApplyScore)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_PostPage, registeredApplyScore)
		}
		c.JSON(200, status)
	}
}

func register(c *gin.Context) {
	date := c.Params.ByName("date")
	collection := c.Params.ByName("collection")
	//		if origin := req.HeaderParameter(restful.HEADER_Origin); origin != "" {
	//			if len(resp.Header().Get(restful.HEADER_AccessControlAllowOrigin)) == 0 {
	//				resp.AddHeader(restful.HEADER_AccessControlAllowOrigin, origin)
	//			}
	//		}

	pp.Println(date)
	l.Output(
		logrus.Fields{
			"Date":       date,
			"Collection": collection,
			"Context":    l.Sprintf(c),
		},
		"Post access to register router",
		l.Debug,
	)

	switch collection {

	case "user":

		userCols := new([]m.UserCol)
		if !c.Bind(userCols) {
			c.JSON(404, gin.H{"status": "wrong queries"})
			return
		}

		status, err := m.RegisterUserColData(*userCols)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_RegisterCol, userCols)
		}
		c.JSON(200, status)

	case "team":

		teamCols := new([]m.TeamCol)
		if !c.Bind(teamCols) {
			c.JSON(404, gin.H{"status": "wrong queries"})
			return
		}

		status, err := m.RegisterTeamColData(date, *teamCols)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_RegisterCol, teamCols)
		}
		c.JSON(200, status)

	case "field":

		fieldCols := new([]m.FieldCol)
		if !c.Bind(fieldCols) {
			c.JSON(404, gin.H{"status": "wrong queries"})
			return
		}

		status, err := m.RegisterFieldColData(date, *fieldCols)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_RegisterCol, fieldCols)
		}
		c.JSON(200, status)

	case "thread":

		requestTakePictureStatus := new(m.RequestTakePictureStatus)
		if !c.Bind(requestTakePictureStatus) {
			c.JSON(404, gin.H{"status": "wrong queries"})
			return
		}

		status, err := m.RegisterThreadImg(requestTakePictureStatus)
		if err != nil {
			l.PutErr(err, l.Trace(), l.E_R_RegisterCol, requestTakePictureStatus)
		}

		c.JSON(200, status)
	}
}
