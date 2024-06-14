package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func DownloadProof(c *gin.Context) {
	var proofRequest struct {
		AssetId   string `json:"asset_id"`
		ProofName string `json:"proof_name"`
	}
	err := c.BindJSON(&proofRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Should Bind JSON with asset_id and proof_name. " + err.Error(),
			Data:    "",
		})
		return
	}
	path, err := services.ValidateAndGetProofFilePath(proofRequest.AssetId, proofRequest.ProofName)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Validate And Get Proof File Path. " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+path)
	c.Header("Content-Disposition", "inline;filename="+path)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.File(path)
	return
}
