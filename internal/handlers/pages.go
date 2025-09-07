package handlers

import (
    "github.com/gin-gonic/gin"
)

// Page handlers (server-side rendered templates)
func IndexPage(c *gin.Context)      { c.HTML(200, "index.html", gin.H{}) }
func AboutPage(c *gin.Context)      { c.HTML(200, "gioi-thieu.html", gin.H{}) }
func APIDocsPage(c *gin.Context)    { c.HTML(200, "api.html", gin.H{}) }

