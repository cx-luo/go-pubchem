// Package router coding=utf-8
// @Project : go-pubchem
// @Time    : 2023/12/12 10:47
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@pharmaron.com
// @File    : router.go
// @Software: GoLand
package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-pubchem/src"
	"go-pubchem/utils"
	"net/http"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter(outputPath string, loglevel string) *gin.Engine {
	// 设置全局 Logger
	logger := utils.SetupLogger(outputPath, loglevel)

	defer logger.Sync()
	// 延迟关闭 logger

	router := gin.New()

	// 使用 Zap 中间件
	router.Use(utils.GinLogger(logger), utils.GinRecovery(logger, true))

	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

var routes = Routes{
	{
		"Swagger",
		http.MethodGet,
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	},

	{
		"GetCidFromSmiles",
		http.MethodPost,
		"/api/pug/getCidFromSmiles",
		src.GetCidFromSmiles,
	},
	{
		"InsertToDbByCid",
		http.MethodPost,
		"/api/insertToDbByCid",
		src.InsertToDbByCid,
	},

	{
		"GetCidFromName",
		http.MethodPost,
		"/api/pug/getCidFromName",
		src.GetCidFromName,
	},
	{
		"GetCmpdWithCasFromCid",
		http.MethodPost,
		"/api/getCmpdWithCasFromCid",
		src.GetCmpdWithCasFromCid,
	},
	{
		"GetCidFromName",
		http.MethodPost,
		"/api/getCidFromName",
		src.GetCidFromName,
	},
}
