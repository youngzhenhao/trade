package servicesrpc

import (
	"context"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"strconv"
	"trade/config"
	"trade/utils"
)

func GetAssetLeaves(ID string, isGroup bool, proofType string) (*universerpc.AssetLeafResponse, error) {
	requset := universerpc.ID{}
	var p universerpc.ProofType
	switch proofType {
	case "issuance":
		p = universerpc.ProofType_PROOF_TYPE_ISSUANCE
	case "transfer":
		p = universerpc.ProofType_PROOF_TYPE_TRANSFER
	default:
		return nil, fmt.Errorf("unknown proof type: %s", proofType)
	}
	requset.ProofType = p

	if isGroup {
		groupId := universerpc.ID_GroupKeyStr{
			GroupKeyStr: ID,
		}
		requset.Id = &groupId
	} else {
		assetId := universerpc.ID_AssetIdStr{
			AssetIdStr: ID,
		}
		requset.Id = &assetId
	}

	leaves, err := getAssetLeaves(&requset)
	if err != nil {
		return nil, err
	}
	return leaves, nil

}
func GetAssetMeta(ID string, isHash bool) (*taprpc.AssetMeta, error) {
	var request taprpc.FetchAssetMetaRequest
	if isHash {
		assetHast := taprpc.FetchAssetMetaRequest_MetaHashStr{
			MetaHashStr: ID,
		}
		request.Asset = &assetHast
	} else {
		assetId := taprpc.FetchAssetMetaRequest_AssetIdStr{
			AssetIdStr: ID,
		}
		request.Asset = &assetId
	}
	assetMeta, err := getAssetMeta(&request)
	if err != nil {
		return nil, err
	}
	return assetMeta, nil
}

func SyncAsset(universe string, id string, isGroupKey bool, proofType string) (*universerpc.SyncResponse, error) {
	request := universerpc.SyncRequest{}
	var p universerpc.ProofType
	switch proofType {
	case "issuance":
		p = universerpc.ProofType_PROOF_TYPE_ISSUANCE
	case "transfer":
		p = universerpc.ProofType_PROOF_TYPE_TRANSFER
	default:
		return nil, fmt.Errorf("unknown proof type: %s", proofType)
	}

	if isGroupKey {
		groupKey := universerpc.ID_GroupKeyStr{
			GroupKeyStr: id,
		}
		request.SyncTargets = append(request.SyncTargets, &universerpc.SyncTarget{
			Id: &universerpc.ID{Id: &groupKey,
				ProofType: p},
		})
	} else {
		assetId := universerpc.ID_AssetIdStr{
			AssetIdStr: id,
		}
		request.SyncTargets = append(request.SyncTargets, &universerpc.SyncTarget{
			Id: &universerpc.ID{Id: &assetId,
				ProofType: p},
		})
	}
	request.UniverseHost = universe
	request.SyncMode = universerpc.UniverseSyncMode_SYNC_ISSUANCE_ONLY
	response, err := syncAsset(&request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getAssetLeaves(request *universerpc.ID) (*universerpc.AssetLeafResponse, error) {
	tapdconf := config.GetConfig().ApiConfig.Tapd

	grpcHost := tapdconf.Host + ":" + strconv.Itoa(tapdconf.Port)
	tlsCertPath := tapdconf.TlsCertPath
	macaroonPath := tapdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := universerpc.NewUniverseClient(conn)
	response, err := client.AssetLeaves(context.Background(), request)
	return response, err
}

func getAssetMeta(request *taprpc.FetchAssetMetaRequest) (*taprpc.AssetMeta, error) {
	tapdconf := config.GetConfig().ApiConfig.Tapd

	grpcHost := tapdconf.Host + ":" + strconv.Itoa(tapdconf.Port)
	tlsCertPath := tapdconf.TlsCertPath
	macaroonPath := tapdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := taprpc.NewTaprootAssetsClient(conn)
	response, err := client.FetchAssetMeta(context.Background(), request)
	return response, err
}

func syncAsset(request *universerpc.SyncRequest) (*universerpc.SyncResponse, error) {
	tapdconf := config.GetConfig().ApiConfig.Tapd

	grpcHost := tapdconf.Host + ":" + strconv.Itoa(tapdconf.Port)
	tlsCertPath := tapdconf.TlsCertPath
	macaroonPath := tapdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := universerpc.NewUniverseClient(conn)
	response, err := client.SyncUniverse(context.Background(), request)
	return response, err
}
