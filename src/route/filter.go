package route

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/emicklei/go-restful"
)

type Config struct {
	Auth struct {
		Admin string `toml:"admin"`
	}
	PagesInfo map[string]PageInfo
}

type PageInfo struct {
	RequireAuth     bool `toml:"requireAuth"`
	ParamaterLength int  `toml:"paramaterLength"`
}

var conf *Config

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	_, err = toml.DecodeFile(path.Join(dir, "../config/config.tml"), &conf)
	if err != nil {
		panic(err)
	}
}

func basicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	param := strings.Split(req.PathParameter("page"), "/")

	if len(param) != conf.PagesInfo[param[0]].ParamaterLength {
		resp.WriteErrorString(
			http.StatusNotFound,
			"404: Page is not found.",
		)
		return
	}
	if conf.PagesInfo[param[0]].RequireAuth {
		encoded := req.Request.Header.Get("Authorization")
		if len(encoded) == 0 || "Basic "+conf.Auth.Admin != encoded {
			resp.AddHeader(
				"WWW-Authenticate",
				"Basic realm=Protected Area",
			)
			resp.WriteErrorString(
				http.StatusUnauthorized,
				"401: Not Authorized",
			)
			return
		}
	}
	chain.ProcessFilter(req, resp)
}
