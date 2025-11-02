package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowInformPage(c *gin.Context) {
	c.HTML(http.StatusOK, "inform.html", gin.H{
		"CurrentPath": c.Request.URL.Path,
	})
}
