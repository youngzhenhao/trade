package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"
	"trade/models"
)

func MakeJsonResult(success bool, error string, data any) string {
	jsr := models.JsonResult{
		Success: success,
		Error:   error,
		Data:    data,
	}
	jsonStr, err := json.Marshal(jsr)
	if err != nil {
		return MakeJsonResult(false, err.Error(), nil)
	}
	return string(jsonStr)
}

func LnMarshalRespString(resp proto.Message) string {
	jsonBytes, err := lnrpc.ProtoJSONMarshalOpts.Marshal(resp)
	if err != nil {
		LogError("unable to decode response", err)
		return ""
	}
	return string(jsonBytes)
}

func TapMarshalRespString(resp proto.Message) string {
	jsonBytes, err := taprpc.ProtoJSONMarshalOpts.Marshal(resp)
	if err != nil {
		LogError("unable to decode response", err)
		return ""
	}
	return string(jsonBytes)
}

func B64DecodeToHex(s string) string {
	byte1, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "DECODE_ERROR"
	}
	return hex.EncodeToString(byte1)
}

type MacaroonCredential struct {
	macaroon string
}

func NewMacaroonCredential(macaroon string) *MacaroonCredential {
	return &MacaroonCredential{macaroon: macaroon}
}

func (c *MacaroonCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"macaroon": c.macaroon}, nil
}

func (c *MacaroonCredential) RequireTransportSecurity() bool {
	return true
}

func GetTimeNow() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func GetTimeSuffixString() string {
	return time.Now().Format("20060102150405")
}

func RoundToDecimalPlace(number float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Round(number*shift) / shift
}

func NewTlsCert(tlsCertPath string) credentials.TransportCredentials {
	cert, err := os.ReadFile(tlsCertPath)
	if err != nil {
		log.Fatalf("Failed to read cert file: %s", err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		log.Fatalf("Failed to append cert")
	}
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    certPool,
	}
	credentialTls := credentials.NewTLS(config)
	return credentialTls
}

func GetMacaroon(macaroonPath string) string {
	macaroonBytes, err := os.ReadFile(macaroonPath)
	if err != nil {
		panic(err)
	}
	macaroon := hex.EncodeToString(macaroonBytes)
	return macaroon
}

func GetConn(grpcHost, tlsCertPath, macaroonPath string) (*grpc.ClientConn, func()) {
	creds := NewTlsCert(tlsCertPath)
	var (
		conn *grpc.ClientConn
		err  error
	)
	if macaroonPath != "" {
		macaroon := GetMacaroon(macaroonPath)
		conn, err = grpc.NewClient(grpcHost, grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(NewMacaroonCredential(macaroon)))
	} else {
		conn, err = grpc.NewClient(grpcHost, grpc.WithTransportCredentials(creds))
	}
	if err != nil {
		LogError("did not connect: grpc.Dial", err)
		return nil, func() {}
	}
	return conn, func() {
		err := conn.Close()
		if err != nil {
			LogError("conn Close Error", err)
		}
	}
}

func GetEnv(key string, filename ...string) string {
	err := godotenv.Load(filename...)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	value := os.Getenv(key)
	return value
}

func ToBTC(sat int) float64 {
	return float64(sat / 1e8)
}

func ToSat(btc float64) int {
	return int(btc * 1e8)
}

func LogInfo(info string) {
	fmt.Printf("%s %s\n", GetTimeNow(), info)
}

func LogInfos(infos ...string) {
	var info string
	for i, _info := range infos {
		if i != 0 {
			info += " "
		}
		info += _info
	}
	fmt.Printf("%s %s\n", GetTimeNow(), info)
}

func LogError(description string, err error) {
	fmt.Printf("%s %s :%v\n", GetTimeNow(), description, err)
}

func CreateFile(filePath, content string) bool {
	dir := filepath.Dir(filePath)
	if dir != "." {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0644)
			if err != nil {
				fmt.Printf("MkdirAll error: %v", err)
				return false
			}
		}
	}
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		LogError("WriteFile error", err)
		return false
	}
	LogInfo("Successes!")
	return true
}

func GetTimestamp() int {
	return int(time.Now().Unix())
}

func Gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return Gcd(b, a%b)
}

func ValueJsonString(value any) string {
	resultJSON, err := json.MarshalIndent(value, "", "\t")
	if err != nil {
		LogError("MarshalIndent error", err)
		return ""
	}
	return string(resultJSON)
}
