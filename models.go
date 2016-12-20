package sal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/thisisfineio/sal/providers/aws/s3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

type Mapping struct {
	AccountID       string
	AccountName     string
	BucketName      string
	ServiceProvider string
}

func (m *Mapping) ProxyManager() (ProxyManager, error) {
	switch m.ServiceProvider {
	case Amazon:
		return &S3Proxy{}, nil
	case Google:
		return nil, errors.New("Google manager unimplemented")
	default:
		return nil, errors.New("Unrecognized ServiceProvider")
	}
}

type Application struct {
	ID             int
	ApiKey         string
	SecretKeyHash  string
	BucketMappings map[string]*Mapping
	Name           string
	Description    string
}

var (
	applications = make(map[string]*Application)
)

var (
	loaders []ApplicationMappingLoader
)

type JSONFileApplicationMapper struct {
	path string
}

func (f *JSONFileApplicationMapper) LoadApplicationMappings() (map[string]*Application, error) {
	data, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	apps := make(map[string]*Application)

	err = json.Unmarshal(data, &apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

type InlineApplicationMapper struct{}

func (f *InlineApplicationMapper) LoadApplicationMappings() (map[string]*Application, error) {
	apps := make(map[string]*Application)
	apps["hsal"] = &Application{
		ID:            1,
		ApiKey:        "THISISNOTAREALKEY",
		SecretKeyHash: "THISISNOTAREALSECRETKEYHASH",
		BucketMappings: map[string]*Mapping{
			"sal.test": {
				AccountID:       "THISISNOTAREALACCOUNTID",
				AccountName:     "THISISNOTAREALACCOUNTNAME",
				BucketName:      "sal",
				ServiceProvider: Amazon,
			},
			"sal-gifs": {
				AccountID:       "THISISNOTAREALACCOUNTID",
				AccountName:     "THISISNOTAREALACCOUNTNAME",
				BucketName:      "sal-gifs",
				ServiceProvider: Amazon,
			},
		},
		Name:        "sal",
		Description: "Storage Abstraction Layer",
	}
	return apps, nil
}

type DatabaseApplicationMapper struct {
	db *sql.DB
}

func (d *DatabaseApplicationMapper) LoadApplicationMappings() (map[string]*Application, error) {
	return nil, errors.New("Unimplemented")
}

func init() {

	db, _ := sql.Open("mysql", "")

	loaders = []ApplicationMappingLoader{
		&DatabaseApplicationMapper{db: db},
		&JSONFileApplicationMapper{path: "test"},
		&InlineApplicationMapper{}}

	for i, loader := range loaders {
		apps, err := loader.LoadApplicationMappings()
		if err != nil {
			if i == len(loaders)-1 {
				log.Fatal(err)
			}
		} else {
			applications = apps
			break
		}
	}
}

type S3Proxy struct{}

func (s *S3Proxy) HandleProxyDownload(mapping *Mapping, c *gin.Context) error {

	b := s3.NewBucket(mapping.BucketName)

	p, err := b.ListPath(strings.TrimLeft(c.Param(v1pathParamString), "/"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return errors.New("S3 Object not fonud")
	} else {
		if len(p.Contents) > 0 {
			// choose the first result
			fi := *p.Contents[0]
			if *fi.Size > DownloadThreshold.Int64() {
				// write to file and then stream the file
				f, err := ioutil.TempFile("", "sal-tmp")
				if err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					return err
				}
				defer os.Remove(f.Name())
				_, err = b.GetObjectAndWriteToWriter(f, *fi.Key)
				if err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					return err
				}
				c.File(f.Name())
			} else {
				// store entirely in memory
				buffer := aws.NewWriteAtBuffer([]byte{})
				_, err := b.GetObjectAndWriteToWriter(buffer, *fi.Key)
				if err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					return err
				}
				readSeeker := bytes.NewReader(buffer.Bytes())
				c.Content(*fi.Key, *fi.LastModified, readSeeker)
			}
		} else {
			c.AbortWithStatus(http.StatusNotFound)
			return errors.New("S3 Object not found")
		}
	}
	return nil
}

func (s *S3Proxy) HandleProxyUpload(mapping *Mapping, c *gin.Context) error {
	//b := s3.NewBucket(app.BucketMappings[c.Param("bucket-name")].BucketName)

	/*tmp, err := ioutil.TempFile("", "hsal-tmp")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}*/
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return err
	}

	spew.Dump(header)
	c.Status(200)
	return nil
}

type GoogleStorageProxy struct{}

func (g *GoogleStorageProxy) HandleProxyDownload(mapping *Mapping, c *gin.Context) error {
	return nil
}

func (g *GoogleStorageProxy) HandleProxyUpload(mapping *Mapping, c *gin.Context) error {
	return nil
}
