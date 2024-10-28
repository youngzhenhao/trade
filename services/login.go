package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
)

const fixedSalt = "bitlongwallet7238baee9c2638664"

// AES密钥（实际应用中应从安全配置获取）
var aesKey = []byte("YourAESKey32BytesLongForSecurity")

func SplitStringAndVerifyChecksum(extstring string) bool {
	originalString, checksum := spilt(extstring)
	if originalString == "" {
		return false
	}
	if checksum == "" {
		return false
	}
	return verifyChecksumWithSalt(originalString, checksum)
}

func spilt(extstring string) (string, string) {
	parts := strings.Split(extstring, "_e_")
	if len(parts) != 2 {
		return "", ""
	}
	originalString := parts[0]
	checksum := parts[1]
	return originalString, checksum
}

func generateMD5WithSalt(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input + fixedSalt))
	return hex.EncodeToString(hasher.Sum(nil))
}

func verifyChecksumWithSalt(originalString, checksum string) bool {
	expectedChecksum := generateMD5WithSalt(originalString)
	return checksum == expectedChecksum
}

func Login(creds *models.User) (string, error) {
	var (
		username = creds.Username
		err      error
	)

	// 检查是否是加密数据
	if isEncrypted(creds.Username) {
		// 解密用户名
		username, err = DecryptAndRestore(creds.Username)
		if err != nil {
			return "", fmt.Errorf("username decryption failed: %v", err)
		}

	}
	//todo 如果手机端都更新到最新代码以下代码需要放开
	//else{
	//	return "", fmt.Errorf("user login failed")
	//}

	var user models.User
	result := middleware.DB.Where("user_name = ?", username).First(&user).Limit(1)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// If there are other database errors, an error is returned
			return "", result.Error
		} else {
			user.Username = username
			password, err := hashPassword(creds.Password)
			if err != nil {
				return "", err
			}
			user.Password = password
			err = btldb.CreateUser(&user)
			if err != nil {
				return "", err
			}
		}
	}
	if !CheckPassword(user.Password, creds.Password) {
		return "", errors.New("invalid credentials")
	}
	token, err := middleware.GenerateToken(username)
	if err != nil {
		return "", err
	}
	creds.Username = username
	return token, nil
}

// isEncrypted 检查数据是否是加密的
func isEncrypted(data string) bool {
	// 检查是否是有效的十六进制字符串
	if _, err := hex.DecodeString(data); err != nil {
		return false
	}

	// 检查长度（AES加密数据的特征）
	if len(data) < 64 {
		return false
	}

	return true
}

// DecryptAndRestore 解密并还原数据
func DecryptAndRestore(encryptedData string) (string, error) {
	if !isEncrypted(encryptedData) {
		return encryptedData, nil // 如果不是加密数据，直接返回
	}

	decrypted, err := aesDecrypt(encryptedData)
	if err != nil {
		return "", err
	}

	restored := removeRandomValues(decrypted)
	return restored, nil
}

// aesDecrypt AES解密
func aesDecrypt(encryptedHex string) (string, error) {
	// 1. 验证输入
	if len(encryptedHex) == 0 {
		return "", fmt.Errorf("empty encrypted data")
	}

	// 2. 解码十六进制
	combined, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", fmt.Errorf("hex decode error: %v", err)
	}

	// 3. 验证长度
	if len(combined) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext size")
	}

	// 4. 分离IV和密文
	iv := combined[:aes.BlockSize]
	ciphertext := combined[aes.BlockSize:]

	// 5. 创建解密器
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	// 6. 创建明文缓冲区
	plaintext := make([]byte, len(ciphertext))

	// 7. 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// 8. 移除填充
	unpadded, err := pkcs7Unpad(plaintext)
	if err != nil {
		return "", err
	}

	return string(unpadded), nil
}

// removeRandomValues 移除随机值
func removeRandomValues(input string) string {
	parts := strings.Split(input, "_")
	var cleaned strings.Builder

	for _, part := range parts {
		if len(part) >= 12 {
			cleaned.WriteString(part[:12])
		} else {
			cleaned.WriteString(part)
		}
	}

	return cleaned.String()
}

// pkcs7Unpad 移除PKCS7填充
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("empty data")
	}

	padding := int(data[length-1])
	if padding > aes.BlockSize || padding == 0 {
		return nil, fmt.Errorf("invalid padding size")
	}

	for i := length - padding; i < length; i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding values")
		}
	}

	return data[:length-padding], nil
}

func ValidateUserAndGenerateToken(creds models.User) (string, error) {
	var (
		username = creds.Username
		err      error
	)

	// 检查是否是加密数据
	if isEncrypted(creds.Username) {
		// 解密用户名
		username, err = DecryptAndRestore(creds.Username)
		if err != nil {
			return "", fmt.Errorf("username decryption failed: %v", err)
		}

	}
	//todo 如果手机端都更新到最新代码以下代码需要放开
	//else{
	//	return "", fmt.Errorf("user login failed")
	//}

	var user models.User
	result := middleware.DB.Where("user_name = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", errors.New("invalid credentials")
	}
	if !CheckPassword(user.Password, creds.Password) {
		originalString, _ := spilt(creds.Password)
		if originalString != "" {
			password, err := hashPassword(originalString)
			if err != nil {
				return "", err
			}
			user.Password = password
			err = btldb.UpdateUser(&user)
			if err != nil {
				return "", err
			}
		}
	}
	token, err := middleware.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (cs *CronService) FiveSecondTask() {
	fmt.Println("5 secs runs")
	log.Println("5 secs runs")
}
