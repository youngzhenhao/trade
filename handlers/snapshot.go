package handlers

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func DownloadSnapshot(c *gin.Context) {
	path := "/root/neutrino/data.zip"
	file, err := os.Open(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File not found"})
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+"data.zip")
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.WriteHeader(http.StatusOK)
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
