package services

import (
	"errors"
	"path"
	"strings"
	"trade/config"
	"trade/utils"
)

// ValidateAndGetProofFilePath
// @dev: Only return nil error when dest is exists
// @notice: Config Tapd's Dir later
func ValidateAndGetProofFilePath(assetId string, proof string) (string, error) {
	var err error
	if strings.Contains(proof, "/") || strings.Contains(proof, "\\") || strings.Contains(proof, "..") {
		err = errors.New("invalid proof, include path")
		return "", err
	}
	if len(assetId) != 64 {
		err = errors.New("wrong assetId length")
		return "", err
	}
	if !utils.IsHexString(assetId) {
		err = errors.New("invalid assetId, not hex")
		return "", err
	}
	proofPath := "data/regtest/proofs"
	dest := path.Join(config.GetLoadConfig().ApiConfig.Tapd.Dir, proofPath, assetId, proof)
	isExist, err := utils.IsPathExists(dest)
	if err != nil {
		return "", err
	}
	if !isExist {
		err = errors.New("proof path does not exist")
		return "", err
	}
	return dest, nil
}
