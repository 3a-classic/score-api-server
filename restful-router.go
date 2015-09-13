package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"./mongo"
	"github.com/emicklei/go-restful"
)

type ProductResource struct {
	// typically reference a DAO (data-access-object)
}

func (p ProductResource) getCol(req *restful.Request, resp *restful.Response) {
	col := req.PathParameter("col")
	log.Println("getting collection data with api:" + col)

	if col == "player" || col == "field" || col == "team" {
		data, err := mongo.GetAllColData(col)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
	}
}

func (p ProductResource) getPage(req *restful.Request, resp *restful.Response) {
	page := req.PathParameter("page")
	team := req.PathParameter("team")
	hole := req.PathParameter("hole")

	log.Println("getting page data with api:" + page)
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
		teamName := team
		holeString := hole
		data, err := mongo.GetScoreEntrySheetPageData(teamName, holeString)
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
	}
}

func (p ProductResource) postOne(req *restful.Request, resp *restful.Response) {
	page := req.PathParameter("page")
	teamName := req.PathParameter("team")
	holeString := req.PathParameter("hole")
	log.Println("post data at " + page)
	if origin := req.HeaderParameter(restful.HEADER_Origin); origin != "" {
		if len(resp.Header().Get(restful.HEADER_AccessControlAllowOrigin)) == 0 {
			resp.AddHeader(restful.HEADER_AccessControlAllowOrigin, origin)
		}
	}
	if page == "scoreEntrySheet" {
		log.Println(req)
		log.Println(resp)

		updatedTeamScore := new(mongo.PostTeamScore)
		err := req.ReadEntity(updatedTeamScore)
		log.Println(updatedTeamScore)
		if err != nil { // bad request
			resp.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}

		status, err := mongo.PostScoreEntrySheetPageData(teamName, holeString, updatedTeamScore)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(status)
		log.Println("updating score team:" + teamName + ", hole: " + holeString)
	}
}

func (p ProductResource) Register(rootPath string) {
	ws := new(restful.WebService)
	ws.Path("/" + rootPath)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/collection/{col)").To(p.getCol).
		Doc("get the product by its id").
		Param(ws.PathParameter("id", "identifier of the collection index").DataType("string")))

	ws.Route(ws.GET("/page/{page)").To(p.getPage).
		Doc("get the page data  by its page").
		Param(ws.PathParameter("page", "identifier of the page index").DataType("string")))

	ws.Route(ws.GET("/page/{page)/{team}").To(p.getPage).
		Doc("get the page data  by its page").
		Param(ws.PathParameter("page", "identifier of the page index").DataType("string")))

	ws.Route(ws.GET("/page/{page)/{team}/{hole}").To(p.getPage).
		Doc("get the page data  by its page").
		Param(ws.PathParameter("page", "identifier of the page index").DataType("string")))

	ws.Route(ws.POST("/page/{page}/{team}/{hole}").To(p.postOne).
		Doc("update or create team score").
		Param(ws.BodyParameter("Product", "a Product (JSON)").DataType("mongo.PostTeamScore")))

	restful.Add(ws)
}

func main() {

	ProductResource{}.Register("api")
	http.ListenAndServe(":8443", nil)
	shutdownHook()
}

func shutdownHook() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	os.Exit(0)
}
