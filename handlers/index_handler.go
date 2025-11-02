package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"CurrentPath": c.Request.URL.Path,
	})
}
