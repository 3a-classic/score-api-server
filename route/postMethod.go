package route

import (
	"log"
	"net/http"

	"../mongo"
	"github.com/emicklei/go-restful"
)

func postOne(req *restful.Request, resp *restful.Response) {
	page := req.PathParameter("page")
	team := req.PathParameter("team")
	hole := req.PathParameter("hole")
	log.Println("post data at " + page)
	if origin := req.HeaderParameter(restful.HEADER_Origin); origin != "" {
		if len(resp.Header().Get(restful.HEADER_AccessControlAllowOrigin)) == 0 {
			resp.AddHeader(restful.HEADER_AccessControlAllowOrigin, origin)
		}
	}
	log.Println(req)
	log.Println(resp)
	switch page {
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
			panic(err)
		}
		resp.WriteAsJson(status)
		log.Println("updating score team:" + team + ", hole: " + hole)
	case "applyScore":
		if hole != "" {
			return
		}
		registeredApplyScore := new(mongo.PostApplyScore)
		err := req.ReadEntity(registeredApplyScore)
		log.Println("registeredApplyScore")
		log.Println(registeredApplyScore)
		if err != nil { // bad request      resp.WriteErrorString(http.StatusBadRequest, err.Error())      return
		}

		status, err := mongo.PostApplyScoreData(team, registeredApplyScore)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(status)
		log.Println("updating apply score:" + team)
	}
}
