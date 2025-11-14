package boot

import (
	"kami/internal/config"
	"kami/internal/handlers"

	// "kami/internal/stmp"
	"kami/internal/utils"

	"kami/internal/logger"
	"kami/internal/middleware"

	"github.com/gin-gonic/gin"
)

func StartRouter() {

	if config.Config.APP_ENV == "production" {

		gin.SetMode(gin.ReleaseMode)
	}
	Router = gin.Default()
	if config.Config.APP_ENV == "production" {

		Router.Use(gin.Recovery())
		// Router.Use(gin.Logger())
	}

	Router.Use(middleware.EnsureIdentity())
	Router.LoadHTMLGlob("web/templates/*")
	Router.Static("/static", "./web/static")

	Router.StaticFile("/robots.txt", "./web/static/robots.txt")
	Router.StaticFile("/sitemap.xml", "./web/static/sitemap.xml")
}
func setupFunctionRoutes() {}
func setupRoutes() {
	Router.GET("/", handlers.IndexPage)
	Router.GET("/gioi-thieu", handlers.AboutPage)
	Router.GET("/about", handlers.AboutPage)
	Router.GET("/api", handlers.APIDocsPage)
	// Router.GET("/debug", func(c *gin.Context){ c.JSON(200, stmp.DATAMAIL) })
	Router.GET("/attachments/:id", handlers.DownloadAttachment)
	Router.GET("/messages", handlers.GetMessagesHandler)
	Router.GET("/messages/:uid", handlers.GetMessageDetailHandler)
	Router.POST("/randomize", handlers.RandomizeHandler)

	// hàm chứa setup route các loại

}
func run() {
	Router.Run(":" + config.Config.PORT)
}

func Boot() {
	// Initialize domain list and blacklist
	_ = utils.LoadDomains("domain.json")
	_ = utils.InitBlacklist("blacklist.json")
	_ = logger.Init("log")

	StartRouter()
	setupFunctionRoutes()
	setupRoutes()
	run()

}
