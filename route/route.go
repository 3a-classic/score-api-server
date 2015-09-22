package route

import "github.com/emicklei/go-restful"

type ProductResource struct {
	// typically reference a DAO (data-access-object)
}

func (p ProductResource) Register(rootPath string) {
	ws := new(restful.WebService)
	ws.Path("/" + rootPath)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)

	// Get URL
	ws.Route(ws.GET("/collection/{col)").To(p.getCol).
		Doc("get the product by its col").
		Param(ws.PathParameter("col", "identifier of the collection index").DataType("string")))

	ws.Route(ws.GET("/page/{page)").To(p.getPage).
		Doc("get the page data  by its page").
		Param(ws.PathParameter("page", "identifier of the page index").DataType("string")))

	ws.Route(ws.GET("/page/{page)/{team}").To(p.getPage).
		Doc("get the page data  by its page").
		Param(ws.PathParameter("page", "identifier of the page index").DataType("string")))

	ws.Route(ws.GET("/page/{page)/{team}/{hole}").To(p.getPage).
		Doc("get the page data  by its page").
		Param(ws.PathParameter("page", "identifier of the page index").DataType("string")))

	//basic auth
	ws.Route(ws.GET("/page/entireScore").Filter(basicAuthenticate).To(p.getSecret))

	//Post URL
	ws.Route(ws.POST("/page/{page}/{team}").To(p.postOne).
		Doc("update apply score").
		Param(ws.BodyParameter("PostApplyScore", "a PostApplyScore  (JSON)").DataType("mongo.PostApplyScore")))

	ws.Route(ws.POST("/page/{page}/{team}/{hole}").To(p.postOne).
		Doc("update or create team score").
		Param(ws.BodyParameter("PostTeamScore", "a PostTeamScore (JSON)").DataType("mongo.PostTeamScore")))

	restful.Add(ws)
}
