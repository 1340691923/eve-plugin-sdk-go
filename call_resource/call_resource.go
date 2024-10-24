package call_resource

import (
	"embed"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/1340691923/eve-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/1340691923/eve-plugin-sdk-go/backend/web_engine"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
)

func NewResourceHandler(webEngine *web_engine.WebEngine, frontendFiles embed.FS) backend.CallResourceHandler {
	//前端页面
	//因为前端所用技术可以进行热更新，所以可进行脱离插件控制
	webEngine.GetGinEngine().Use(ServeFrontendFiles("/", EmbedFolder(frontendFiles, "dist")))

	return httpadapter.New(webEngine.Handler())
}

func ServeFrontendFiles(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}
