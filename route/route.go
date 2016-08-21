package route

import (
	"github.com/gin-gonic/gin"
)

func Register() {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/collection/:col", getCol)
		v1.GET("/page/*page", getPage)

		v1.POST("/page/:page", postOne)
		v1.POST("/page/:page/:team", postOne)
		v1.POST("/page/:page/:team/:hole", postOne)
		v1.POST("/register/:collection", register)
		v1.POST("/register/:collection/:date", register)
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/collection/:col", getCol)
		v2.GET("/page/*page", getPage)

		v2.POST("/page/:page", postOne)
		v2.POST("/page/:page/:team", postOne)
		v2.POST("/page/:page/:team/:hole", postOne)
		v2.POST("/register/:collection", register)
	}

	r.Run(":8080")

}
