package utils

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
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

func MakeJsonErrorResult(code models.ErrCode, errorString string, data any) string {
	jsr := models.JsonResult{
		Error: errorString,
		Code:  code,
		Data:  data,
	}
	if errors.Is(code, models.SUCCESS) {
		jsr.Success = true
	} else {
		jsr.Success = false
	}
	jstr, err := json.Marshal(jsr)
	if err != nil {
		return MakeJsonErrorResult(models.DefaultErr, err.Error(), nil)
	}
	return string(jstr)
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
		log.Fatalf("Failed to read cert file: %s(%s)", err, tlsCertPath)
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

func AppendErrorInfo(err error, info string) error {
	if err == nil {
		err = errors.New("[nil err]")
	}
	return errors.New(err.Error() + "; " + info)
}

func AppendError(e error) func(error) error {
	return func(err error) error {
		if e == nil {
			e = errors.New("")
		}
		if err == nil {
			return e
		}
		if e.Error() == "" {
			e = err
			return e
		}
		e = errors.New(e.Error() + "; " + err.Error())
		return e
	}
}

func AppendInfo(s string) func(string) string {
	return func(info string) string {
		if info == "" {
			return s
		}
		if s == "" {
			s = info
			return s
		}
		s = s + "; " + info
		return s
	}
}

func InfoAppendError(i string) func(error) error {
	e := errors.New(i)
	return func(err error) error {
		if err == nil {
			return e
		}
		if e.Error() == "" {
			e = err
			return e
		}
		e = errors.New(e.Error() + "; " + err.Error())
		return e
	}
}

func ErrorAppendInfo(e error) func(string) error {
	return func(info string) error {
		if e == nil {
			e = errors.New("")
		}
		if info == "" {
			return e
		}
		if e.Error() == "" {
			e = errors.New(info)
			return e
		}
		info = e.Error() + "; " + info
		e = errors.New(info)
		return e
	}
}

func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsHexString(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

func SwapValue[T any](a *T, b *T) {
	temp := *a
	*a = *b
	*b = temp
}

func SwapInt(a *int, b *int) {
	*a ^= *b
	*b ^= *a
	*a ^= *b
}

func ToLowerWords(s string) string {
	var result strings.Builder
	for i, char := range s {
		if i > 0 && char >= 'A' && char <= 'Z' {
			temp := result.String()
			if len(temp) > 0 && temp[len(temp)-1] != ' ' {
				result.WriteRune(' ')
			}
		}
		result.WriteRune(char)
	}
	return strings.ToLower(result.String())
}

func ToLowerWordsWithHyphens(s string) string {
	var result strings.Builder
	for i, char := range s {
		if char == ' ' {
			continue
		}
		if i > 0 && char >= 'A' && char <= 'Z' {
			temp := result.String()
			if len(temp) > 0 && temp[len(temp)-1] != ' ' {
				result.WriteRune('-')
			}
		}
		result.WriteRune(char)
	}
	return strings.ToLower(result.String())
}

func ToLowerWordsWithUnderscores(s string) string {
	var result strings.Builder
	for i, char := range s {
		if char == ' ' {
			continue
		}
		if i > 0 && char >= 'A' && char <= 'Z' {
			temp := result.String()
			if len(temp) > 0 && temp[len(temp)-1] != ' ' {
				result.WriteRune('_')
			}
		}
		result.WriteRune(char)
	}
	return strings.ToLower(result.String())
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func ToCamelWord(s string, isByUnderline bool, isLowerCaseInitial bool) string {
	var sli []string
	if isByUnderline {
		sli = strings.Split(s, "_")
	} else {
		sli = strings.Split(s, " ")
	}
	var result strings.Builder
	for _, word := range sli {
		if result.String() == "" && isLowerCaseInitial {
			result.WriteString(word)
		} else {
			result.WriteString(FirstUpper(word))
		}
	}
	return result.String()
}

func CreateTestMainFile(testPath string, testFuncName string) {
	dirPath := path.Join(testPath, ToLowerWordsWithHyphens(testFuncName))
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	filePath := path.Join(dirPath, "main.go")
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(f)
	content := []byte("package main\n\nfunc main() {\n\n}\n")
	_, err = f.Write(content)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(filePath, "has been created successfully!")
}

func BuildTestMainFile(testPath string, testFuncName string) {
	if strings.HasPrefix(testFuncName, ".\\") && strings.HasSuffix(testFuncName, ".exe") {
		testFuncName, _ = strings.CutPrefix(testFuncName, ".\\")
		testFuncName, _ = strings.CutSuffix(testFuncName, ".exe")
	}
	dirPath := path.Join(testPath, ToLowerWordsWithHyphens(testFuncName))
	filePath := path.Join(dirPath, "main.go")
	executableFileName := testFuncName + ".exe"
	cmd := exec.Command("go", "build", "-o", executableFileName, filePath)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(path.Join(testPath, executableFileName), "has been built successfully!")
}

func GetFunctionName(i any) string {
	completeName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	s := strings.Split(completeName, ".")
	return s[len(s)-1]
}

// GetTransactionAndIndexByOutpoint
// @dev: Split outpoint
func GetTransactionAndIndexByOutpoint(outpoint string) (transaction string, index string) {
	result := strings.Split(outpoint, ":")
	return result[0], result[1]
}

func GetTxidFromOutpoint(outpoint string) (string, error) {
	txid, indexStr := GetTransactionAndIndexByOutpoint(outpoint)
	if txid == "" || indexStr == "" {
		return "", errors.New("txid or index is empty")
	}
	return txid, nil
}

func OutpointToTransactionAndIndex(outpoint string) (transaction string, index string) {
	result := strings.Split(outpoint, ":")
	return result[0], result[1]
}

// GetFilesNameInPathWithoutDir
// @Description: Get name slice of files in specific path
func GetFilesNameInPathWithoutDir(path string) ([]string, error) {
	var fileSlice []string
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			fileSlice = append(fileSlice, file.Name())
		}
	}
	return fileSlice, nil
}

// GetFilesNameInPathRecursive
// @Description: Get name slice of files in specific path and recursive folders
func GetFilesNameInPathRecursive(path string) ([]string, error) {
	var fileSlice []string
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			var pathFiles []string
			// @dev: Recursive fetch
			pathFiles, err = GetFilesNameInPathRecursive(filepath.Join(path, file.Name()))
			if err != nil {
				return nil, err
			}
			fileSlice = append(fileSlice, pathFiles...)
		} else {
			fileSlice = append(fileSlice, filepath.Join(path, file.Name()))
		}
	}
	return fileSlice, nil
}

func RemoveStringInNameForFileSlice(fileSlice []string, s string) ([]string, error) {
	var newSlice []string
	for _, file := range fileSlice {
		name := strings.ReplaceAll(file, s, "")
		err := os.Rename(file, name)
		if err != nil {
			return nil, err
		}
		newSlice = append(newSlice, name)
	}
	return newSlice, nil
}

// GetRandomNumber
// @Description: Return a random number whose range is (0,maxValue).
func GetRandomNumber(maxValue int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var randNumber int
	for randNumber == 0 {
		randNumber = rand.Intn(maxValue)
	}
	return randNumber
}

// GetRandomNumberSlice
// @Description: Use this to generate a slice of random number
func GetRandomNumberSlice(maxValue int, n int) ([]int, error) {
	// @dev: If the number of times the same random number is obtained exceeds this value,
	// add this random number to the slice of random number results,
	// so that duplicates can be obtained relatively randomly
	// when the range of random numbers is small
	// (the number of random numbers that can be obtained is less than
	// the number of random numbers that need to be generated).
	// The third time when a random number duplicates, it will be accepted.
	const retryTimes = 2
	if maxValue < 1 || n < 0 {
		return nil, errors.New("max value or slice length is negative")
	}
	var slice []int
	randNumberMapGeneratedTimes := make(map[int]int)
	for len(slice) < n {
		// @dev: Return a random number whose range is (0,maxValue].
		randNumber := GetRandomNumber(maxValue + 1)
		randNumberMapGeneratedTimes[randNumber]++
		if randNumberMapGeneratedTimes[randNumber] == 1 || randNumberMapGeneratedTimes[randNumber] > 1+retryTimes {
			slice = append(slice, randNumber)
			randNumberMapGeneratedTimes[randNumber] = 1
		} else {
			// @dev: length of slice will not add
			continue
		}
	}
	return slice, nil
}

func GetNowTimeStringWithHyphens() string {
	now := time.Now().Format("2006-01-02-15-04-05.000000")
	now = strings.ReplaceAll(now, ".", "-")
	return now
}

func GetNewFileLogger(file *os.File, prefix string) *log.Logger {
	return log.New(file, prefix, log.Lshortfile|log.Ldate|log.Ltime)
}

func GetLogFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
}

func WriteToLogFile(filename string, prefix string, content string) error {
	file, err := GetLogFile(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	logger := GetNewFileLogger(file, prefix)
	logger.Println("\n", content)
	return nil
}

func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func DateTimeStringToTime(dateTime string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		err = AppendErrorInfo(err, "LoadLocation")
		return time.Time{}, err
	}
	return time.ParseInLocation(time.DateTime, dateTime, loc)
}

func DateTimeStringToTimeWithFormat(dateTime string, format string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		err = AppendErrorInfo(err, "LoadLocation")
		return time.Time{}, err
	}
	return time.ParseInLocation(format, dateTime, loc)
}

func GetTimePeriodString(_time time.Time) string {
	hour := _time.Hour()
	var period string
	switch {
	case hour >= 5 && hour < 8:
		period = "清晨"
	case hour >= 8 && hour < 12:
		period = "上午"
	case hour >= 12 && hour < 14:
		period = "中午"
	case hour >= 14 && hour < 17:
		period = "下午"
	case hour >= 17 && hour < 19:
		period = "傍晚"
	default:
		period = "夜间"
	}
	return period
}

func TimeFormatCN(_time time.Time) string {
	format := _time.Format("2006年01月02日 15:04:05")
	split := strings.Split(format, " ")
	format = split[0] + " " + GetTimePeriodString(_time) + " " + split[1]
	return format
}

func RandString(r *rand.Rand, length int) string {
	var dic = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = dic[r.Intn(len(dic))]
	}
	return string(b)
}

func Sha256(data any) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(jsonData)
	hashString := fmt.Sprintf("%x", hash)
	return hashString, nil
}
