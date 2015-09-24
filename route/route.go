package route

import "github.com/emicklei/go-restful"

func Register() {
	ws := new(restful.WebService)
	ws.
		Path("/api").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// Get URL
	ws.Route(ws.GET("/collection/{col)").
		To(getCol).
		Doc("get the product by its col").
		Param(ws.PathParameter("col", "identifier of the collection index").DataType("string")))

	ws.Route(ws.GET("/page/{page:*}").
		Filter(basicAuthenticate).
		To(getPage).
		Doc("get the page data  by its page").
		Param(ws.PathParameter("page", "identifier of the page index").DataType("string")))

	//Post URL
	ws.Route(ws.POST("/page/{page}/{team}").
		To(postOne).
		Doc("update apply score").
		Param(ws.BodyParameter("PostApplyScore", "a PostApplyScore  (JSON)").DataType("mongo.PostApplyScore")))

	ws.Route(ws.POST("/page/{page}/{team}/{hole}").
		To(postOne).
		Doc("update or create team score").
		Param(ws.BodyParameter("PostTeamScore", "a PostTeamScore (JSON)").DataType("mongo.PostTeamScore")))

	restful.Add(ws)
}
