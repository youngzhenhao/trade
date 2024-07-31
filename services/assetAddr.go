package services

import (
	"errors"
	"time"
	"trade/models"
)

func GetAssetAddrsByUserId(userId int) (*[]models.AssetAddr, error) {
	return ReadAssetAddrsByUserId(userId)
}

func ProcessAssetAddrSetRequest(userId int, username string, assetAddrSetRequest *models.AssetAddrSetRequest) *models.AssetAddr {
	var assetAddr models.AssetAddr
	assetAddr = models.AssetAddr{
		Encoded:          assetAddrSetRequest.Encoded,
		AssetId:          assetAddrSetRequest.AssetId,
		AssetType:        assetAddrSetRequest.AssetType,
		Amount:           assetAddrSetRequest.Amount,
		GroupKey:         assetAddrSetRequest.GroupKey,
		ScriptKey:        assetAddrSetRequest.ScriptKey,
		InternalKey:      assetAddrSetRequest.InternalKey,
		TapscriptSibling: assetAddrSetRequest.TapscriptSibling,
		TaprootOutputKey: assetAddrSetRequest.TaprootOutputKey,
		ProofCourierAddr: assetAddrSetRequest.ProofCourierAddr,
		AssetVersion:     assetAddrSetRequest.AssetVersion,
		DeviceID:         assetAddrSetRequest.DeviceID,
		UserId:           userId,
		Username:         username,
	}
	return &assetAddr
}

func IsAssetAddrChanged(assetAddrByAddrEncoded *models.AssetAddr, old *models.AssetAddr) bool {
	if assetAddrByAddrEncoded == nil || old == nil {
		return true
	}
	if assetAddrByAddrEncoded.Encoded != old.Encoded {
		return true
	}
	if assetAddrByAddrEncoded.AssetId != old.AssetId {
		return true
	}
	if assetAddrByAddrEncoded.AssetType != old.AssetType {
		return true
	}
	if assetAddrByAddrEncoded.Amount != old.Amount {
		return true
	}
	if assetAddrByAddrEncoded.GroupKey != old.GroupKey {
		return true
	}
	if assetAddrByAddrEncoded.ScriptKey != old.ScriptKey {
		return true
	}
	if assetAddrByAddrEncoded.InternalKey != old.InternalKey {
		return true
	}
	if assetAddrByAddrEncoded.TapscriptSibling != old.TapscriptSibling {
		return true
	}
	if assetAddrByAddrEncoded.TaprootOutputKey != old.TaprootOutputKey {
		return true
	}
	if assetAddrByAddrEncoded.ProofCourierAddr != old.ProofCourierAddr {
		return true
	}
	if assetAddrByAddrEncoded.AssetVersion != old.AssetVersion {
		return true
	}
	if assetAddrByAddrEncoded.DeviceID != old.DeviceID {
		return true
	}
	if assetAddrByAddrEncoded.UserId != old.UserId {
		return true
	}
	return false
}

func CheckAssetAddrIfUpdate(assetAddr *models.AssetAddr) (*models.AssetAddr, error) {
	if assetAddr == nil {
		return nil, errors.New("nil asset addr")
	}
	assetAddrByAddrEncoded, err := ReadAssetAddrByAddrEncoded(assetAddr.Encoded)
	if err != nil {
		return assetAddr, nil
	}
	if !IsAssetAddrChanged(assetAddrByAddrEncoded, assetAddr) {
		return assetAddrByAddrEncoded, nil
	}
	assetAddrByAddrEncoded.Encoded = assetAddr.Encoded
	assetAddrByAddrEncoded.AssetId = assetAddr.AssetId
	assetAddrByAddrEncoded.AssetType = assetAddr.AssetType
	assetAddrByAddrEncoded.Amount = assetAddr.Amount
	assetAddrByAddrEncoded.GroupKey = assetAddr.GroupKey
	assetAddrByAddrEncoded.ScriptKey = assetAddr.ScriptKey
	assetAddrByAddrEncoded.InternalKey = assetAddr.InternalKey
	assetAddrByAddrEncoded.TapscriptSibling = assetAddr.TapscriptSibling
	assetAddrByAddrEncoded.TaprootOutputKey = assetAddr.TaprootOutputKey
	assetAddrByAddrEncoded.ProofCourierAddr = assetAddr.ProofCourierAddr
	assetAddrByAddrEncoded.AssetVersion = assetAddr.AssetVersion
	assetAddrByAddrEncoded.DeviceID = assetAddr.DeviceID
	assetAddrByAddrEncoded.UserId = assetAddr.UserId
	return assetAddrByAddrEncoded, nil
}

func CreateOrUpdateAssetAddr(addr *models.AssetAddr) (err error) {
	var assetAddr *models.AssetAddr
	assetAddr, err = CheckAssetAddrIfUpdate(addr)
	return UpdateAssetAddr(assetAddr)
}

func GetAssetAddrsByScriptKey(scriptKey string) (*[]models.AssetAddr, error) {
	return ReadAssetAddrsByScriptKey(scriptKey)
}

func GetAssetAddrsByEncoded(encoded string) (*models.AssetAddr, error) {
	return ReadAssetAddrByAddrEncoded(encoded)
}

func GetAllAssetAddrs() (*[]models.AssetAddr, error) {
	return ReadAllAssetAddrs()
}

func UpdateUsernameByUserIdAll() error {
	allAssetAddrs, err := GetAllAssetAddrs()
	if allAssetAddrs == nil || *allAssetAddrs == nil || len(*allAssetAddrs) == 0 {
		return nil
	}
	if err != nil {
		return err
	}
	for i, assetAddr := range *allAssetAddrs {
		var name string
		name, err = IdToName(assetAddr.UserId)
		if err != nil {
			continue
		}
		(*allAssetAddrs)[i].Username = name
	}
	return UpdateAssetAddrs(allAssetAddrs)
}

type AssetAddrSimplified struct {
	UpdatedAt time.Time `json:"updated_at"`
	Encoded   string    `json:"encoded"`
	AssetId   string    `json:"asset_id"`
	Amount    int       `json:"amount"`
	ScriptKey string    `json:"script_key"`
	DeviceID  string    `json:"device_id"`
	Username  string    `json:"username"`
}

func AssetAddrToAssetAddrSimplified(assetAddr models.AssetAddr) AssetAddrSimplified {
	return AssetAddrSimplified{
		UpdatedAt: assetAddr.UpdatedAt,
		Encoded:   assetAddr.Encoded,
		AssetId:   assetAddr.AssetId,
		Amount:    assetAddr.Amount,
		ScriptKey: assetAddr.ScriptKey,
		DeviceID:  assetAddr.DeviceID,
		Username:  assetAddr.Username,
	}
}

func AssetAddrSliceToAssetAddrSimplifiedSlice(assetAddrs *[]models.AssetAddr) *[]AssetAddrSimplified {
	var assetAddrSimplified []AssetAddrSimplified
	for _, assetAddr := range *assetAddrs {
		assetAddrSimplified = append(assetAddrSimplified, AssetAddrToAssetAddrSimplified(assetAddr))
	}
	return &assetAddrSimplified
}

func GetAllAssetAddrSimplified() (*[]AssetAddrSimplified, error) {
	allAssetAddrs, err := ReadAllAssetAddrs()
	if err != nil {
		return nil, err
	}
	return AssetAddrSliceToAssetAddrSimplifiedSlice(allAssetAddrs), nil
}

type UsernameAssetAddr struct {
	Username   string                 `json:"username"`
	AssetAddrs *[]AssetAddrSimplified `json:"asset_addrs"`
}

type AssetIdAssetAddr struct {
	AssetId    string                 `json:"asset_id"`
	AssetAddrs *[]AssetAddrSimplified `json:"asset_addrs"`
}

type UsernameAndAssetIdAssetAddr struct {
	Username          string              `json:"username"`
	AssetIdAssetAddrs *[]AssetIdAssetAddr `json:"asset_id_asset_addrs"`
}

func AssetAddrSimplifiedSliceToUsernameMapAssetAddrs(assetAddrSimplified *[]AssetAddrSimplified) *map[string]*[]AssetAddrSimplified {
	if assetAddrSimplified == nil {
		return nil
	}
	usernameMapAssetAddrs := make(map[string]*[]AssetAddrSimplified)
	for _, assetAddr := range *assetAddrSimplified {
		assetAddrs, ok := usernameMapAssetAddrs[assetAddr.Username]
		if !ok {
			usernameMapAssetAddrs[assetAddr.Username] = &[]AssetAddrSimplified{assetAddr}
		} else {
			*assetAddrs = append(*assetAddrs, assetAddr)
		}
	}
	return &usernameMapAssetAddrs
}

func UsernameMapAssetAddrsToUsernameAssetAddrs(usernameMapAssetAddrs *map[string]*[]AssetAddrSimplified) *[]UsernameAssetAddr {
	if usernameMapAssetAddrs == nil {
		return nil
	}
	var usernameAssetAddrs []UsernameAssetAddr
	for username, assetAddrs := range *usernameMapAssetAddrs {
		usernameAssetAddrs = append(usernameAssetAddrs, UsernameAssetAddr{
			Username:   username,
			AssetAddrs: assetAddrs,
		})
	}
	return &usernameAssetAddrs
}

func AssetAddrSimplifiedSliceToUsernameAssetAddrs(assetAddrSimplified *[]AssetAddrSimplified) *[]UsernameAssetAddr {
	if assetAddrSimplified == nil {
		return nil
	}
	usernameMapAssetAddrs := AssetAddrSimplifiedSliceToUsernameMapAssetAddrs(assetAddrSimplified)
	usernameAssetAddrs := UsernameMapAssetAddrsToUsernameAssetAddrs(usernameMapAssetAddrs)
	return usernameAssetAddrs
}

func AssetAddrSimplifiedSliceToAssetIdMapAssetAddrs(assetAddrSimplified *[]AssetAddrSimplified) *map[string]*[]AssetAddrSimplified {
	if assetAddrSimplified == nil {
		return nil
	}
	assetIdMapAssetAddrs := make(map[string]*[]AssetAddrSimplified)
	for _, assetAddr := range *assetAddrSimplified {
		assetAddrs, ok := assetIdMapAssetAddrs[assetAddr.AssetId]
		if !ok {
			assetIdMapAssetAddrs[assetAddr.AssetId] = &[]AssetAddrSimplified{assetAddr}
		} else {
			*assetAddrs = append(*assetAddrs, assetAddr)
		}
	}
	return &assetIdMapAssetAddrs
}

func AssetIdMapAssetAddrsToAssetIdAssetAddrs(assetIdMapAssetAddrs *map[string]*[]AssetAddrSimplified) *[]AssetIdAssetAddr {
	if assetIdMapAssetAddrs == nil {
		return nil
	}
	var assetIdAssetAddrs []AssetIdAssetAddr
	for assetId, assetAddrs := range *assetIdMapAssetAddrs {
		assetIdAssetAddrs = append(assetIdAssetAddrs, AssetIdAssetAddr{
			AssetId:    assetId,
			AssetAddrs: assetAddrs,
		})
	}
	return &assetIdAssetAddrs
}

func AssetAddrSimplifiedSliceToAssetIdAssetAddrs(assetAddrSimplified *[]AssetAddrSimplified) *[]AssetIdAssetAddr {
	if assetAddrSimplified == nil {
		return nil
	}
	assetIdMapAssetAddrs := AssetAddrSimplifiedSliceToAssetIdMapAssetAddrs(assetAddrSimplified)
	assetIdAssetAddrs := AssetIdMapAssetAddrsToAssetIdAssetAddrs(assetIdMapAssetAddrs)
	return assetIdAssetAddrs
}

func AssetAddrSimplifiedSliceToUsernameAndAssetIdAssetAddrs(assetAddrSimplified *[]AssetAddrSimplified) *[]UsernameAndAssetIdAssetAddr {
	if assetAddrSimplified == nil {
		return nil
	}
	var usernameAndAssetIdAssetAddrs []UsernameAndAssetIdAssetAddr
	usernameAssetAddrs := AssetAddrSimplifiedSliceToUsernameAssetAddrs(assetAddrSimplified)
	for _, usernameAssetAddr := range *usernameAssetAddrs {
		assetIdAssetAddrs := AssetAddrSimplifiedSliceToAssetIdAssetAddrs(usernameAssetAddr.AssetAddrs)
		usernameAndAssetIdAssetAddrs = append(usernameAndAssetIdAssetAddrs, UsernameAndAssetIdAssetAddr{
			Username:          usernameAssetAddr.Username,
			AssetIdAssetAddrs: assetIdAssetAddrs,
		})
	}
	return &usernameAndAssetIdAssetAddrs
}

func GetAllUsernameAndAssetIdAssetAddrs() (*[]UsernameAndAssetIdAssetAddr, error) {
	allAssetAddrSimplified, err := GetAllAssetAddrSimplified()
	if err != nil {
		return nil, err
	}
	usernameAndAssetIdAssetAddrs := AssetAddrSimplifiedSliceToUsernameAndAssetIdAssetAddrs(allAssetAddrSimplified)
	return usernameAndAssetIdAssetAddrs, nil
}
