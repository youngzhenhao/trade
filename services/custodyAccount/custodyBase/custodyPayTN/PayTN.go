package custodyPayTN

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

type PayToNpubKey struct {
	NpubKey string
	AssetId string
	Amount  float64
	Time    int64
	Vision  uint8
}

func (p *PayToNpubKey) Encode() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	hexData := hex.EncodeToString(data)
	return "ptn" + hexData, nil
}
func (p *PayToNpubKey) Decode(encoded string) error {
	if !strings.HasPrefix(encoded, "ptn") {
		return errors.New("无效的编码字符串: 缺少前缀 'ptn'")
	}

	hexData := encoded[3:] // 去掉前缀 "PTN"
	data, err := hex.DecodeString(hexData)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, p)
}

func HashEncodedString(encoded string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(encoded))
	if err != nil {
		return "", err
	}
	hashedBytes := hash.Sum(nil)
	return hex.EncodeToString(hashedBytes), nil
}
