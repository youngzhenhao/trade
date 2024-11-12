package services

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"strings"
	"sync"
	"time"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
)

func GenerateNonce() (string, error) {
	// 创建一个 16 字节的随机值
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	// 使用随机值和时间戳生成哈希
	hash := sha256.New()
	hash.Write(randomBytes)
	hash.Write([]byte(time.Now().String()))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

var nonceStore = make(map[string]time.Time)

// 保存 nonce 并设置过期时间
func StoreNonceInRedis(username string, tokenString string) (string, error) {
	// 从 Redis 中检查是否已有 nonce
	userName, err := middleware.RedisGet(tokenString)
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}
	if userName != "" {
		// 如果已存在 nonce，直接返回现有的 token
		return tokenString, nil
	}

	// 如果 Redis 操作出错并且不是 key 不存在的错误
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	// 设置 nonce 的有效期（默认为 10 分钟）
	redisSetTimeMinute := 10 // 可以从配置中获取
	expiration := time.Duration(redisSetTimeMinute) * time.Minute

	// 将 username 与 nonce 的映射存储在 Redis 中
	err = middleware.RedisSet(tokenString, username+"_nonce", expiration)
	if err != nil {
		return "", err
	}
	// 返回生成的 noncexuy
	return tokenString, nil
}

// 验证 nonce 是否存在并有效
func VerifyNonce(nonce string, usernameRef string) bool {
	// 从 Redis 获取 nonce，如果不存在则返回 false
	username, err := middleware.RedisGet(nonce)
	if err != nil || username == "" {
		return false // nonce 不存在或 Redis 出错
	}
	if usernameRef+"_nonce" != username {
		return false
	}
	// 验证成功后删除 nonce，确保只能使用一次
	if err := middleware.RedisDel(nonce); err != nil {
		return true // 删除失败，可能需要处理错误
	}
	return true
}

// DeviceIDGenerator 设备ID生成器
type DeviceIDGenerator struct {
	lastTimestamp int64
	sequence      int64
	mutex         sync.Mutex
}

// NewDeviceIDGenerator 创建新的设备ID生成器实例
func NewDeviceIDGenerator() *DeviceIDGenerator {
	return &DeviceIDGenerator{
		lastTimestamp: 0,
		sequence:      0,
		mutex:         sync.Mutex{},
	}
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func (g *DeviceIDGenerator) GenerateDeviceID(prefix string, randomLength int) (string, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// 获取当前时间戳（毫秒）
	currentTimestamp := time.Now().UnixNano() / 1e6

	// 如果在同一毫秒内，序列号自增
	if currentTimestamp == g.lastTimestamp {
		g.sequence++
	} else {
		g.sequence = 0
		g.lastTimestamp = currentTimestamp
	}

	// 生成随机字符串
	randomStr, err := GenerateRandomString(randomLength)
	if err != nil {
		return "", err
	}

	// 构建设备ID
	var builder strings.Builder

	// 添加前缀（如果有）
	if prefix != "" {
		builder.WriteString(prefix)
		builder.WriteString("-")
	}
	// 添加时间戳
	builder.WriteString(fmt.Sprintf("%d", currentTimestamp))

	// 添加序列号（防止同一毫秒内的冲突）
	builder.WriteString(fmt.Sprintf("%03d", g.sequence))

	// 添加随机字符串
	builder.WriteString("-")
	builder.WriteString(randomStr)
	return builder.String(), nil
}

func GetDeviceID() (string, error) {
	generator := NewDeviceIDGenerator()
	deviceID, err := generator.GenerateDeviceID("DEV", 8)
	if err != nil {
		return "", err
	}
	return deviceID, nil
}
func generateKeyAndSalt(password []byte) ([]byte, []byte, error) {
	// 生成一个随机的盐值
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, nil, err
	}
	// 使用PBKDF2生成派生密钥
	key := pbkdf2.Key(password, salt, 10000, 32, sha256.New)

	return key, salt, nil
}
func encrypt(plainText, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 使用PKCS7填充模式
	blockSize := block.BlockSize()
	padding := blockSize - len(plainText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	plainText = append(plainText, padText...)

	iv := make([]byte, blockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(plainText))
	mode.CryptBlocks(encrypted, plainText)

	// 将IV和加密后的密文拼接在一起，以便解密时使用
	result := append(iv, encrypted...)

	return base64.StdEncoding.EncodeToString(result), nil
}
func BuildEncrypt(deviceID string) (string, string, error) {
	password := []byte("thisisaverysecretkey1234567890")
	// 生成密钥和盐值
	key, salt, err := generateKeyAndSalt(password)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	encryptedID, err := encrypt([]byte(deviceID), key)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(salt), encryptedID, nil
}
func decrypt(cipherText, key []byte) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(cipherText))
	if err != nil {
		return "", err
	}

	// 提取IV
	iv := decoded[:16]
	encrypted := decoded[16:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	// 去除填充
	unpadding := int(decrypted[len(decrypted)-1])
	return string(decrypted[:len(decrypted)-unpadding]), nil
}
func BuildDecrypt(saltBase64 string, encryptedDeviceID string) string {
	password := []byte("thisisaverysecretkey1234567890")
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		fmt.Println("解码盐值失败:", err)
		return ""
	}
	key := pbkdf2.Key(password, salt, 10000, 32, sha256.New)
	decryptedID, err := decrypt([]byte(encryptedDeviceID), key)
	if err != nil {
		fmt.Println("解密失败:", err)
		return ""
	}
	return decryptedID
}

func checkNpublicExists(npublic string) (bool, string) {
	device, err := btldb.ReadDeviceManagerByNpubKey(npublic)
	if err != nil {
		return false, ""
	}
	return true, device.DeviceID
}
func ProcessDeviceRequest(nonce, nPubKey string) (string, string, error) {
	// Step 0: Validate input parameters
	if nonce == "" || nPubKey == "" {
		return "", "", errors.New("nonce or nPubKey cannot be empty")
	}

	// Step 1: Validate nonce
	if !VerifyNonce(nonce, nPubKey) {
		return "", "", errors.New("invalid or expired nonce")
	}

	// Step 2: Check if Npublic exists in the database using ReadDeviceManagerByNpubKey
	flag, deviceId := checkNpublicExists(nPubKey)
	if !flag {
		var device models.DeviceManager
		deviceID, err := GetDeviceID()
		if err != nil {
			return "", "", err
		}
		device.DeviceID = deviceID
		device.Status = 1
		device.NpubKey = nPubKey
		encryptDeviceID, encodedSalt, err := BuildEncrypt(deviceID)
		if err != nil {
			return "", "", err
		}
		device.EncryptDeviceID = encodedSalt
		err = btldb.CreateDeviceManager(&device)
		if err != nil {
			return "", "", err
		}
		return encryptDeviceID, encodedSalt, nil
	}

	// Step 3: Encrypt deviceId with salt
	deviceID := deviceId // Assume DeviceID is a field in your models.DeviceManager struct
	encryptDeviceID, encodedSalt, err := BuildEncrypt(deviceID)
	if err != nil {
		return "", "", err
	}
	// Step 4: Return encrypted deviceId and salt
	return encryptDeviceID, encodedSalt, nil
}
