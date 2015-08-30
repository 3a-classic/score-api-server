package main

import (
	"./mongo"
	//	"fmt"
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

func (p ProductResource) getOne(req *restful.Request, resp *restful.Response) {
	col := req.PathParameter("col")
	//	id := "1"
	log.Println("getting data with api:" + col)
	if col == "player" {
		data, err := collection.GetAllPlayerJson()
		//	fmt.Println(data)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
	} else if col == "field" {
		data, err := collection.GetAllFieldJson()
		//	fmt.Println(data)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
	} else if col == "team" {
		data, err := collection.GetAllTeamJson()
		//	fmt.Println(data)
		if err != nil {
			panic(err)
		}
		resp.WriteAsJson(data)
		//	resp.WriteEntity(Product{Id: id, Title: "test"})
		//	resp.WriteEntity(data)
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
	//	ws.Consumes(restful.MIME_XML)
	//	ws.Produces(restful.MIME_XML)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)

	//	ws.Route(ws.GET("/{id}").To(p.getOne).
	ws.Route(ws.GET("/{col)").To(p.getOne).
		Doc("get the product by its id").
		Param(ws.PathParameter("id", "identifier of the product").DataType("string")))

	ws.Route(ws.POST("").To(p.postOne).
		Doc("update or create a product").
		Param(ws.BodyParameter("Product", "a Product (XML)").DataType("main.Product")))

	restful.Add(ws)
}

func main() {
	ProductResource{}.Register("api")
	http.ListenAndServe(":8443", nil)
}
