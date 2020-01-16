package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// 签名加密
func Sign(sysParam map[string]string, busParam string, method string, key string) (string, error) {
	if len(sysParam) == 0 || sysParam == nil {
		return "", errors.New("sysParam Valid")
	}
	if method == "HmacSHA256" {
		// 生成加密参数
		decodeHexKey, _ := DecodeHexUpper(key)
		busContent := EncodeAES256HexUpper(busParam, decodeHexKey)
		return signWithSHA256(sysParam, busContent, key), nil
	} else if method == "RSAWithMD5" {
		return signWithRSA(sysParam, busParam, key), nil
	} else {
		return "", errors.New("method   Valid")
	}
}

// sha256方法加密
func signWithSHA256(sysParam map[string]string, busParam string, key string) string {
	if len(busParam) > 0 && busParam != "" && len(strings.TrimSpace(busParam)) > 0 {
		sysParam["content"] = busParam
	}
	var keys []string
	for k := range sysParam {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	buf := make([]string, 200)
	buf = append(buf, key)
	for _, v := range keys {
		if !strings.EqualFold("sign", v) {
			sysVal := v + sysParam[v]
			buf = append(buf, sysVal)
		}
	}
	buf = append(buf, key)
	newString := ""
	for _, v := range buf {
		newString += fmt.Sprintf("%s", v)
	}
	newKey, _ := DecodeHexUpper(key)
	retStr := encodeHmacSHA256HexUpper(newString, newKey)
	return retStr
}

// 中介方法
func encodeHmacSHA256HexUpper(data string, key []byte) string {
	dataByte := []byte(data)
	encodeHmac := encodeHmacSHA256(dataByte, key)
	retStr := bytesToHexString(encodeHmac)
	return strings.ToUpper(retStr)
}

// rsa方式加密
func signWithRSA(sysParam map[string]string, busParam string, key string) string {
	// 暂时用不到,不做整理了
	fmt.Println(sysParam, busParam, key)
	return ""
}

// 16进制字符串转换成byte
func DecodeHexUpper(str string) ([]byte, error) {
	return hex.DecodeString(strings.ToLower(str))
}

// 中介方法
func EncodeAES256HexUpper(data string, key []byte) string {
	dataByte := []byte(data)
	newByte, _ := AesECBEncrypt(dataByte, key)
	retStr := encodeHexUpper([]byte(newByte))
	return retStr
}

// boss返回结果解密
func DecodeAES256HexUpper(data string, key []byte) string {
	newData := strings.ToLower(data)
	dataByte, _ := hex.DecodeString(newData)
	newByte, _ := AesECBDecrypt(dataByte, key)
	return string(newByte)
}

// 16进制转换字符串-结果大写
func encodeHexUpper(data []byte) string {
	str := bytesToHexString(data)
	return strings.ToUpper(str)
}

// byte转16进制字符串
func bytesToHexString(b []byte) string {
	return hex.EncodeToString(b)
}

// 16进制字符串转bytes
func hexStringToBytes(hexString string) []byte {
	newHexString := strings.ToUpper(hexString)
	length := len(hexString) / 2
	newByte := []byte(newHexString)
	retByte := make([]byte, length)
	for i := 0; i < length; i++ {
		pos := i * 2
		retByte[i] = byte(byteToByte(newByte[pos])<<4 | byteToByte(newByte[pos+1]))
	}
	return retByte
}

// byte转换
func byteToByte(b byte) int {
	byteList := []byte("0123456789ABCDEF")
	var ret int
	for k, v := range byteList {
		if string(v) == string(b) {
			ret = k
		}
	}
	return ret
}

// 获取md5加密字符串
func getMD5Str(str string) string {
	md5Data := md5.New()
	md5Data.Reset()
	md5Data.Write([]byte(str))
	retString := bytesToHexString(md5Data.Sum(nil))
	return retString
}

// Hmac-sha256加密
func encodeHmacSHA256(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// md5加密
func MD5Util(s string) string {
	hexDigits := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	byteStr := []byte(s)
	MD5 := md5.New()
	MD5.Write(byteStr)
	MD5Data := MD5.Sum([]byte(nil))
	NewByte := make([]byte, len(MD5Data)*2)
	k := 0
	for i := 0; i < len(MD5Data); i++ {
		byte0 := MD5Data[i]
		NewByte[k] = hexDigits[byte0>>4&15]
		k++
		NewByte[k] = hexDigits[byte0&15]
		k++
	}
	return string(NewByte)
}

// AES/CBC解密数据--不加填充,数据加密
func AesCBCEncrypt(source string, key string) (string, error) {
	// 生成16进制加密key
	newKey := hexStringToBytes(getMD5Str(key))
	block, err := aes.NewCipher(newKey)
	if err != nil {
		return "", err
	}
	// 数据处理
	dataLen := len([]byte(source))
	m := dataLen % 16
	if m != 0 {
		for i := 0; i < 16-m; i++ {
			source = source + " "
		}
	}
	newByte := []byte(source)
	// 初始向量IV必须是唯一
	iv := hexStringToBytes(getMD5Str(key))
	// block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)
	encryptData := make([]byte, len(newByte))
	mode.CryptBlocks(encryptData, newByte)
	return bytesToHexString(encryptData), nil
}

// AES/CBC解密数据--不加填充,数据解密
func AesCBCDecrypt(source, key string) (string, error) {
	// 生成16进制加密key
	newKey := hexStringToBytes(getMD5Str(key))
	block, err := aes.NewCipher(newKey)
	if err != nil {
		return "", err
	}
	// 16进制转换
	decodeBytes := hexStringToBytes(source)
	iv := hexStringToBytes(getMD5Str(key))
	mode := cipher.NewCBCDecrypter(block, iv)
	retData := make([]byte, len(decodeBytes))
	mode.CryptBlocks(retData, decodeBytes)
	return string(retData), nil
}

// PKCS7加填充/和PKCS5填充一样,只是填充字段多少的区别
func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// PKCS7解填充/和PKCS5填充一样,只是填充字段多少的区别
func PKCS7UnPadding(encrypt []byte) []byte {
	length := len(encrypt)
	unPadding := int(encrypt[length-1])
	return encrypt[:(length - unPadding)]
}

// AES/ECB/PKCS7模式加密--签名加密方式
func AesECBEncrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	ecb := NewECBEncryptEr(block)
	// 加PKCS7填充
	content := PKCS7Padding(data, block.BlockSize())
	encryptData := make([]byte, len(content))
	// 生成加密数据
	ecb.CryptBlocks(encryptData, content)
	return encryptData, nil
}

// AES/ECB/PKCS7模式解密--签名解密方式
func AesECBDecrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	ecb := NewECBDecryptEr(block)
	retData := make([]byte, len(data))
	ecb.CryptBlocks(retData, data)
	// 解PKCS7填充
	retData = PKCS7UnPadding(retData)
	return retData, nil
}
