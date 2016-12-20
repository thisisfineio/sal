package sal

import (
	"log"
	"net/http"

	"github.com/alistanis/size"

	"github.com/gin-gonic/gin"
)

const (
	v1applicationParamString = "application-name"
	v1bucketParamString      = "bucket-name"
	v1pathParamString        = "path"
	V1PathMappingString      = "/:" + v1applicationParamString + "/:" + v1bucketParamString + "/*" + v1pathParamString
)

var (
	DownloadThreshold = size.GigaBytes(1)
)

type V1Handler struct{}

func (v *V1Handler) HandleGet(c *gin.Context) {
	app, ok := applications[c.Param(v1applicationParamString)]
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	mapping, ok := app.BucketMappings[c.Param(v1bucketParamString)]
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	cm, err := mapping.ProxyManager()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = cm.HandleProxyDownload(mapping, c)
	if err != nil {
		log.Println(err)
	}
}

func (v *V1Handler) HandlePost(c *gin.Context) {
	app, ok := applications[c.Param(v1applicationParamString)]
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	mapping, ok := app.BucketMappings[c.Param(v1bucketParamString)]
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	cm, err := mapping.ProxyManager()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = cm.HandleProxyUpload(mapping, c)
	if err != nil {
		log.Println(err)
	}
}

func HandleAuthorization(c *gin.Context) error {
	return nil
}
