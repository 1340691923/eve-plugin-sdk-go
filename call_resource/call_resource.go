// call_resource包提供资源调用和HTTP服务功能
package call_resource

// 导入所需的包
import (
	// 导入嵌入文件系统包
	"embed"
	// 导入backend核心功能包
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	// 导入HTTP适配器包
	"github.com/1340691923/eve-plugin-sdk-go/backend/resource/httpadapter"
	// 导入Web引擎包
	"github.com/1340691923/eve-plugin-sdk-go/backend/web_engine"
	// 导入Gin框架
	"github.com/gin-gonic/gin"
	// 导入文件系统接口包
	"io/fs"
	// 导入HTTP包
	"net/http"
	// 导入性能分析包
	"net/http/pprof"
)

// NewResourceHandler 创建一个新的资源处理器
func NewResourceHandler(webEngine *web_engine.WebEngine, frontendFiles embed.FS) backend.CallResourceHandler {

	// 配置前端文件服务
	// 因为前端所用技术可以进行热更新，所以可进行脱离插件控制
	webEngine.GetGinEngine().Use(ServeFrontendFiles("/", EmbedFolder(frontendFiles, "dist")))

	// 附加性能分析路由
	attachPprof(webEngine.GetGinEngine())

	// 返回HTTP适配器实例
	return httpadapter.New(webEngine.Handler())
}

// attachPprof 将pprof调试工具附加到Gin路由器上
func attachPprof(router *gin.Engine) {
	// 创建pprof路由组
	pprofGroup := router.Group("/debug/pprof")
	{
		// 注册各个pprof处理器
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapH(http.HandlerFunc(pprof.Handler("allocs").ServeHTTP)))
		pprofGroup.GET("/block", gin.WrapH(http.HandlerFunc(pprof.Handler("block").ServeHTTP)))
		pprofGroup.GET("/goroutine", gin.WrapH(http.HandlerFunc(pprof.Handler("goroutine").ServeHTTP)))
		pprofGroup.GET("/heap", gin.WrapH(http.HandlerFunc(pprof.Handler("heap").ServeHTTP)))
		pprofGroup.GET("/mutex", gin.WrapH(http.HandlerFunc(pprof.Handler("mutex").ServeHTTP)))
		pprofGroup.GET("/threadcreate", gin.WrapH(http.HandlerFunc(pprof.Handler("threadcreate").ServeHTTP)))
	}
}

// ServeFrontendFiles 创建一个用于服务前端文件的中间件
func ServeFrontendFiles(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	// 创建文件服务器
	fileserver := http.FileServer(fs)
	// 如果有URL前缀，则从请求路径中剥离该前缀
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	// 返回处理函数
	return func(c *gin.Context) {
		// 如果文件存在，则使用文件服务器处理请求
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			// 终止后续中间件处理
			c.Abort()
		}
	}
}

// EmbedFolder 将嵌入式文件系统的子目录转换为ServeFileSystem接口
func EmbedFolder(fsEmbed embed.FS, targetPath string) ServeFileSystem {
	// 从嵌入式文件系统获取指定路径的子文件系统
	fsys, err := fs.Sub(fsEmbed, targetPath)
	// 如果出错，则触发panic
	if err != nil {
		panic(err)
	}
	// 返回包装后的嵌入式文件系统
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

// ServeFileSystem 定义文件服务的接口
type ServeFileSystem interface {
	// 继承http.FileSystem接口
	http.FileSystem
	// 检查路径是否存在
	Exists(prefix string, path string) bool
}

// embedFileSystem 实现ServeFileSystem接口的嵌入式文件系统
type embedFileSystem struct {
	// 内嵌http.FileSystem
	http.FileSystem
}

// Exists 检查给定路径的文件是否存在
func (e embedFileSystem) Exists(prefix string, path string) bool {
	// 尝试打开文件
	_, err := e.Open(path)
	// 如果出错，则文件不存在
	if err != nil {
		return false
	}
	// 文件存在
	return true
}
