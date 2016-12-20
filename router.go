package sal

import "github.com/gin-gonic/gin"

const (
	v1Prefix = "/v1"
)

func Run() error {

	router := gin.Default()

	v1group := router.Group(v1Prefix)
	{
		v1 := &V1Handler{}
		v1group.GET(V1PathMappingString, v1.HandleGet)
		v1group.POST(V1PathMappingString, v1.HandlePost)

	}

	return router.Run()
}
