package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"trade/models"
	"trade/services"
	"trade/services/assetsyncinfo"
)

func DownloadProof(c *gin.Context) {
	AssetId := c.Param("asset_id")
	ProofName := c.Param("proof_name")
	path, err := services.ValidateAndGetProofFilePath(AssetId, ProofName)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Validate And Get Proof File Path. " + err.Error(),
			Data:    nil,
			Code:    models.DefaultErr,
		})
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+ProofName)
	c.Header("Content-Disposition", "inline;filename="+ProofName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.File(path)
	return
}

func DownloadProof2(c *gin.Context) {
	AssetId := c.Param("asset_id")
	ProofName := c.Param("proof_name")
	path, err := services.ValidateAndGetProofFilePath(AssetId, ProofName)
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
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+ProofName)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.WriteHeader(http.StatusOK)
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// SycnAssetInfo 拉取资产同步信息
func SycnAssetInfo(c *gin.Context) {
	req := assetsyncinfo.SyncInfoRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if req.Id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("request error")})
		return
	}
	assetSyncInfo, err := assetsyncinfo.GetAssetSyncInfo(&req)
	if err != nil && errors.Is(err, assetsyncinfo.AssetNotFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Errorf("server error")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accountModel": assetSyncInfo})
}
