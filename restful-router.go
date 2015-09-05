package main

import (
	"./mongo"
	//	"fmt"
	"github.com/emicklei/go-restful"
	//	"io"
	"log"
	"net/http"
	//	"reflect"
	//	"strconv"
)

// This example shows how to use methods as RouteFunctions for WebServices.
// The ProductResource has a Register() method that creates and initializes
// a WebService to expose its methods as REST operations.
// The WebService is added to the restful.DefaultContainer.
// A ProductResource is typically created using some data access object.
//
// GET http://localhost:8080/products/1
// POST http://localhost:8080/products
// <Product><Id>1</Id><Title>The First</Title></Product>

//type PostTeamScore struct {
//	Member []string
//	Stroke []int
//	putt   []int
//	excnt  int
//}

type ProductResource struct {
	// typically reference a DAO (data-access-object)
}

func (p ProductResource) getCol(req *restful.Request, resp *restful.Response) {
	col := req.PathParameter("col")
	log.Println("getting collection data with api:" + col)
	if col == "player" {
		data, err := mongo.GetAllPlayerCol()
		if err != nil {
			panic(err)
		}
		//		fmt.Println(data)
		//		fmt.Println(reflect.ValueOf(data).Type())
		//		fmt.Println(fmt.Sprint("%+v", data))
		resp.WriteAsJson(data)
	} else if col == "field" {
		data, err := mongo.GetAllFieldCol()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
	} else if col == "team" {
		data, err := mongo.GetAllTeamCol()
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
	if page == "index" {
		if team != "" || hole != "" {
			return
		}
		data, err := mongo.GetIndexPageData()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
	} else if page == "leadersBoard" {
		if team != "" || hole != "" {
			return
		}
		data, err := mongo.GetLeadersBoardPageData()
		if err != nil {
			panic(err)
		}
		//		fmt.Println(data == nil)
		resp.WriteAsJson(data)
	} else if page == "scoreEntrySheet" {
		//		log.Println(reflect.ValueOf(team).Type())
		//		log.Println(reflect.ValueOf(hole).Type())
		teamName := team
		holeString := hole
		data, err := mongo.GetScoreEntrySheetPageData(teamName, holeString)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
	} else if page == "scoreViewSheet" {
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
	//	var chain *restful.FilterChain
	if origin := req.HeaderParameter(restful.HEADER_Origin); origin != "" {
		// prevent duplicate header
		if len(resp.Header().Get(restful.HEADER_AccessControlAllowOrigin)) == 0 {
			resp.AddHeader(restful.HEADER_AccessControlAllowOrigin, origin)
		}
	}
	if page == "scoreEntrySheet" {
		log.Println("before chain")
		log.Println(req)
		log.Println(resp)
		//		chain.ProcessFilter(req, resp)

		updatedTeamScore := new(mongo.PostTeamScore)
		err := req.ReadEntity(updatedTeamScore)
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

	//	ws.Route(ws.GET("/{id}").To(p.getOne).
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
		//	ws.Route(ws.POST("").To(p.postOne).
		Doc("update or create team score").
		Param(ws.BodyParameter("Product", "a Product (JSON)").DataType("mongo.PostTeamScore")))

	restful.Add(ws)
}

func main() {

	ProductResource{}.Register("api")
	http.ListenAndServe(":8443", nil)
}
