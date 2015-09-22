package route

import (
	"log"

	"../mongo"
	"github.com/emicklei/go-restful"
)

func (p ProductResource) getSecret(req *restful.Request, resp *restful.Response) {

	log.Println("getting page data with api: entireScore")
	data, err := mongo.GetEntireScorePageData()
	if err != nil {
		panic(err)
	}
	resp.WriteAsJson(data)

}

func basicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	encoded := req.Request.Header.Get("Authorization")
	// usr/pwd = admin/admin
	// real code does some decoding
	//  if len(encoded) == 0 || "Basic YWRtaW46YWRtaW4=" != encoded {

	if len(encoded) == 0 || "Basic M2E6Y2xhc3NpYw==" != encoded {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	}
	chain.ProcessFilter(req, resp)
}
