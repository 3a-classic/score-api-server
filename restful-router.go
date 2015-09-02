package main

import (
	"./mongo"
	"fmt"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
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

type Product struct {
	Id, Title string
}

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
		fmt.Println(data)
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
	log.Println("getting page data with api:" + page)
	if page == "index" {
		data, err := mongo.GetIndexPageData()
		if err != nil {
			panic(err)
		}
		fmt.Println(data)
		resp.WriteAsJson(data)
	} else if page == "leaderBoard" {
		data, err := mongo.GetLeaderBoardPageData()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
	} else if page == "scoreEntrySheet" {
		data, err := mongo.GetScoreEntrySheetPageData()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
	} else if page == "scoreViewSheet" {
		data, err := mongo.GetScoreViewSheetPageData()
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
	}
}

func (p ProductResource) postOne(req *restful.Request, resp *restful.Response) {
	updatedProduct := new(Product)
	err := req.ReadEntity(updatedProduct)
	if err != nil { // bad request
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	log.Println("updating product with id:" + updatedProduct.Id)
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

	ws.Route(ws.POST("").To(p.postOne).
		Doc("update or create a product").
		Param(ws.BodyParameter("Product", "a Product (XML)").DataType("main.Product")))

	restful.Add(ws)
}

func main() {
	ProductResource{}.Register("api")
	http.ListenAndServe(":8443", nil)
}
