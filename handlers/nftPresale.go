package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"trade/btlLog"
	"trade/models"
	"trade/services"
)

// @dev: Get

func GetNftPresaleByAssetId(c *gin.Context) {
	username := c.MustGet("username").(string)
	_, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	assetId := c.Query("asset_id")
	nftPresale, err := services.GetNftPresaleByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetNftPresalesByAssetIdErr,
			Data:    nil,
		})
		return
	}
	noMetaStr := c.Query("no_meta")
	noMeta, err := strconv.ParseBool(noMetaStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	noWhitelistStr := c.Query("no_whitelist")
	noWhitelist, err := strconv.ParseBool(noWhitelistStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	result := services.NftPresaleToNftPresaleSimplified(nftPresale, noMeta, noWhitelist)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    result,
	})
}

func GetNftPresaleByBatchGroupId(c *gin.Context) {
	username := c.MustGet("username").(string)
	_, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	batchGroupIdStr := c.Query("batch_group_id")
	batchGroupId, err := strconv.Atoi(batchGroupIdStr)
	if err != nil {
		btlLog.PreSale.Error("Atoi err:%v", err)
	}
	nftPresale, err := services.GetNftPresaleByBatchGroupId(batchGroupId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetNftPresaleByBatchGroupIdErr,
			Data:    nil,
		})
		return
	}
	noMetaStr := c.Query("no_meta")
	noMeta, err := strconv.ParseBool(noMetaStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	noWhitelistStr := c.Query("no_whitelist")
	noWhitelist, err := strconv.ParseBool(noWhitelistStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	result := services.NftPresaleSliceToNftPresaleSimplifiedSlice(nftPresale, noMeta, noWhitelist)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    result,
	})
}

func GetLaunchedNftPresale(c *gin.Context) {
	username := c.MustGet("username").(string)
	_, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	nftPresales, err := services.GetLaunchedNftPresales()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetLaunchedNftPresalesErr,
			Data:    nil,
		})
		return
	}
	noMetaStr := c.Query("no_meta")
	noMeta, err := strconv.ParseBool(noMetaStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	noWhitelistStr := c.Query("no_whitelist")
	noWhitelist, err := strconv.ParseBool(noWhitelistStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	result := services.NftPresaleSliceToNftPresaleSimplifiedSlice(nftPresales, noMeta, noWhitelist)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    result,
	})
}

func GetUserBoughtNftPresale(c *gin.Context) {
	username := c.MustGet("username").(string)
	userId, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	nftPresales, err := services.GetNftPresalesByBuyerUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetNftPresalesByBuyerUserIdErr,
			Data:    nil,
		})
		return
	}
	noMetaStr := c.Query("no_meta")
	noMeta, err := strconv.ParseBool(noMetaStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	noWhitelistStr := c.Query("no_whitelist")
	noWhitelist, err := strconv.ParseBool(noWhitelistStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	result := services.NftPresaleSliceToNftPresaleSimplifiedSlice(nftPresales, noMeta, noWhitelist)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    result,
	})
}

func GetNftPresaleByGroupKeyPurchasable(c *gin.Context) {
	username := c.MustGet("username").(string)
	_, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	groupKey := c.Query("group_key")
	nftPresales, err := services.GetNftPresaleByGroupKeyPurchasable(groupKey)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetNftPresaleByGroupKeyErr,
			Data:    nil,
		})
		return
	}
	noMetaStr := c.Query("no_meta")
	noMeta, err := strconv.ParseBool(noMetaStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	noWhitelistStr := c.Query("no_whitelist")
	noWhitelist, err := strconv.ParseBool(noWhitelistStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	result := services.NftPresaleSliceToNftPresaleSimplifiedSlice(nftPresales, noMeta, noWhitelist)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    result,
	})
}

func GetNftPresaleNoGroupKeyPurchasable(c *gin.Context) {
	username := c.MustGet("username").(string)
	_, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	nftPresales, err := services.GetNftPresaleNoGroupKeyPurchasable()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetNftPresaleByGroupKeyErr,
			Data:    nil,
		})
		return
	}
	noMetaStr := c.Query("no_meta")
	noMeta, err := strconv.ParseBool(noMetaStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	noWhitelistStr := c.Query("no_whitelist")
	noWhitelist, err := strconv.ParseBool(noWhitelistStr)
	if err != nil {
		btlLog.PreSale.Error("ParseBool err:%v", err)
	}
	result := services.NftPresaleSliceToNftPresaleSimplifiedSlice(nftPresales, noMeta, noWhitelist)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    result,
	})
}

// @dev: Purchase

func BuyNftPresale(c *gin.Context) {
	username := c.MustGet("username").(string)
	userId, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	var buyNftPresaleRequest models.BuyNftPresaleRequest
	err = c.ShouldBindJSON(&buyNftPresaleRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	err = services.BuyNftPresale(userId, username, buyNftPresaleRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.BuyNftPresaleErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    nil,
	})
}

// @dev: Query

func QueryNftPresaleGroupKeyPurchasable(c *gin.Context) {
	username := c.MustGet("username").(string)
	_, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	groupKeys, err := services.GetAllNftPresaleGroupKeyPurchasable()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetLaunchedNftPresalesErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    groupKeys,
	})
}

// @dev: Set

func SetNftPresale(c *gin.Context) {
	var nftPresaleSetRequest models.NftPresaleSetRequest
	err := c.ShouldBindJSON(&nftPresaleSetRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	nftPresale := services.ProcessNftPresale(&nftPresaleSetRequest)
	// @dev: Store AssetMeta
	{
		assetId := nftPresaleSetRequest.AssetId
		err = services.StoreAssetMetaIfNotExist(assetId)
		if err != nil {
			btlLog.PreSale.Error("api StoreAssetMetaIfNotExist err:%v", err)
		}
	}
	err = services.CreateNftPresale(nftPresale)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.CreateNftPresaleErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    nil,
	})
}

func SetNftPresales(c *gin.Context) {
	var nftPresaleSetRequests []models.NftPresaleSetRequest
	err := c.ShouldBindJSON(&nftPresaleSetRequests)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	nftPresales := services.ProcessNftPresales(&nftPresaleSetRequests)
	// @dev: Store AssetMetas
	{
		var assetIds []string
		for _, nftPresaleSetRequest := range nftPresaleSetRequests {
			assetIds = append(assetIds, nftPresaleSetRequest.AssetId)
		}
		err = services.StoreAssetMetasIfNotExist(assetIds)
		if err != nil {
			btlLog.PreSale.Error("api StoreAssetMetasIfNotExist err:%v", err)
		}
	}
	err = services.CreateNftPresales(nftPresales)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.CreateNftPresalesErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    nil,
	})
}

func ReSetFailOrCanceledNftPresale(c *gin.Context) {
	err := services.ReSetFailOrCanceledNftPresale()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ReSetFailOrCanceledNftPresaleErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    nil,
	})
}

// @dev: launch batch group

func LaunchNftPresaleBatchGroup(c *gin.Context) {
	var nftPresaleBatchGroupLaunchRequest models.NftPresaleBatchGroupLaunchRequest
	err := c.ShouldBindJSON(&nftPresaleBatchGroupLaunchRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	// @dev: Store AssetMetas
	{
		var assetIds []string
		for _, nftPresaleSetRequest := range *(nftPresaleBatchGroupLaunchRequest.NftPresaleSetRequests) {
			assetIds = append(assetIds, nftPresaleSetRequest.AssetId)
		}
		err = services.StoreAssetMetasIfNotExist(assetIds)
		if err != nil {
			btlLog.PreSale.Error("api StoreAssetMetasIfNotExist err:%v", err)
		}
	}
	// @dev: Process and create db records
	err = services.ProcessNftPresaleBatchGroupLaunchRequestAndCreate(&nftPresaleBatchGroupLaunchRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ProcessNftPresaleBatchGroupLaunchRequestAndCreateErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    nil,
	})
}

func QueryNftPresaleBatchGroup(c *gin.Context) {
	username := c.MustGet("username").(string)
	_, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	stateStr := c.Query("state")
	state, err := strconv.Atoi(stateStr)
	if err != nil {
		btlLog.PreSale.Error("Atoi err:%v", err)
	}
	batchGroups, err := services.GetBatchGroups(models.NftPresaleBatchGroupState(state))
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetBatchGroupsErr,
			Data:    nil,
		})
		return
	}
	result := services.NftPresaleBatchGroupSliceToNftPresaleBatchGroupSimplifiedSlice(batchGroups)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    result,
	})
}
