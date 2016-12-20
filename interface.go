package sal

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	Amazon = "Amazon"
	Google = "Google"
)

type Object interface {
	Name() string
	Bucket() string
	Data() []byte
	ModTime() time.Time
	Provider() string
}

type ApplicationMappingLoader interface {
	LoadApplicationMappings() (map[string]*Application, error)
}

type ProxyManager interface {
	HandleProxyUpload(*Mapping, *gin.Context) error
	HandleProxyDownload(*Mapping, *gin.Context) error
}
