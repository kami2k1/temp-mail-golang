package handlers

import (
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"

    "kami/internal/config"
    "kami/internal/repository"
    "kami/internal/stmp"
    "kami/internal/utils"
)


func DownloadAttachment(c *gin.Context) {
    id := c.Param("id")
    att, ok := stmp.GetAttachment(id)
    if !ok {
        c.JSON(404, gin.H{"error": "attachment not found"})
        return
    }
    if att.MIME != "" { c.Header("Content-Type", att.MIME) }
    if att.Filename != "" { c.Header("Content-Disposition", "attachment; filename=\""+att.Filename+"\"") }
  
    c.Header("Cache-Control", "public, max-age=31536000, immutable")
    c.Header("Content-Length", fmt.Sprintf("%d", len(att.Data)))
    c.Data(200, att.MIME, att.Data)
}


func GetMessagesHandler(c *gin.Context) {
    stmp.TouchActivity()
    email := c.GetString("email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "no identity"})
        return
    }
    resp := repository.GetMessages(email)
    c.JSON(http.StatusOK, resp)
}


func GetMessageDetailHandler(c *gin.Context) {
    stmp.TouchActivity()
    email := c.GetString("email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "no identity"})
        return
    }
    uidStr := c.Param("uid")
    var uid uint32
    if _, err := fmt.Sscanf(uidStr, "%d", &uid); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uid"})
        return
    }
    msg, ok := repository.GetMessageByUID(email, uid)
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, msg)
}


func RandomizeHandler(c *gin.Context) {
    stmp.TouchActivity()
    email, ok := utils.GenerateRandomEmail()
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate identity"})
        return
    }
    tok, exp, _ := utils.GenerateToken(config.Config.JWT_SECRET, email, 24*time.Hour)
    maxAge := int(time.Until(exp).Seconds())
    c.SetCookie("token", tok, maxAge, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{"email": email})
}
