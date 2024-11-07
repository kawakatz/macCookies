package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"database/sql"
	"errors"
	"os"

	"github.com/kawakatz/macCookies/pkg/types"
	"github.com/kawakatz/macCookies/pkg/utils"

	"github.com/cixtor/binarycookies"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/pbkdf2"
)

var (
	queryFirefoxCookie  = `SELECT name, value, host, path, creationTime, expiry, isSecure, isHttpOnly FROM moz_cookies`
	queryChromiumCookie = `SELECT name, encrypted_value, host_key, path, creation_utc, expires_utc, is_secure, is_httponly, has_expires, is_persistent FROM cookies`
)

var (
	errSecurityKeyIsEmpty = errors.New("input [security find-generic-password -wa 'Chrome'] in terminal")
	//errPasswordIsEmpty    = errors.New("password is empty")
	errDecryptFailed = errors.New("decrypt encrypt value failed")
	//errDecodeASN1Failed   = errors.New("decode ASN1 data failed")
	//errEncryptedLength = errors.New("length of encrypted password less than block size")
)

func FirefoxCookies(cookiesFile string) []types.Cookie {
	cookiesDB, _ := sql.Open("sqlite3", "file:"+cookiesFile)
	rows, _ := cookiesDB.Query(queryFirefoxCookie)

	var decryptedCookies []types.Cookie
	for rows.Next() {
		var (
			name, value, host, path string
			isSecure, isHttpOnly    int
			creationTime, expiry    int64
		)

		_ = rows.Scan(&name, &value, &host, &path, &creationTime, &expiry, &isSecure, &isHttpOnly)

		decryptedCookies = append(decryptedCookies, types.Cookie{
			KeyName:    name,
			Host:       host,
			Path:       path,
			IsSecure:   utils.IntToBool(isSecure),
			IsHTTPOnly: utils.IntToBool(isHttpOnly),
			CreateDate: utils.TimeStampFormat(creationTime / 1000000),
			ExpireDate: utils.TimeStampFormat(expiry),
			Value:      value,
		})
	}
	return decryptedCookies
}

func SafariCookies(cookiesFile string) []types.Cookie {
	f, _ := os.Open(cookiesFile)
	cook := binarycookies.New(f)
	pages, _ := cook.Decode()

	var decryptedCookies []types.Cookie
	for _, page := range pages {
		for _, cookie := range page.Cookies {
			decryptedCookies = append(decryptedCookies, types.Cookie{
				KeyName:    string(cookie.Name),
				Host:       string(cookie.Domain),
				Path:       string(cookie.Path),
				IsSecure:   cookie.Secure,
				IsHTTPOnly: cookie.HttpOnly,
				CreateDate: utils.TimeStampFormat(cookie.Creation.Unix() / 1000000),
				ExpireDate: cookie.Expires,
				Value:      string(cookie.Value),
			})
		}
	}

	return decryptedCookies
}

func MacPassword2SecretKey(chromePassword string) []byte {
	chromeSecret := []byte(chromePassword)
	chromeSalt := []byte("saltysalt")
	secretKey := pbkdf2.Key(chromeSecret, chromeSalt, 1003, 16, sha1.New)

	return secretKey
}

func ChromeCookies(cookiesFile string, secretKey []byte, osType string) []types.Cookie {
	cookiesDB, _ := sql.Open("sqlite3", "file:"+cookiesFile)
	rows, _ := cookiesDB.Query(queryChromiumCookie)

	var decryptedCookies []types.Cookie
	for rows.Next() {
		var (
			key, host, path                               string
			isSecure, isHTTPOnly, hasExpire, isPersistent int
			createDate, expireDate                        int64
			value, encryptValue                           []byte
		)

		_ = rows.Scan(&key, &encryptValue, &host, &path, &createDate, &expireDate, &isSecure, &isHTTPOnly, &hasExpire, &isPersistent)

		cookie := types.Cookie{
			KeyName:      key,
			Host:         host,
			Path:         path,
			EncryptValue: encryptValue,
			IsSecure:     utils.IntToBool(isSecure),
			IsHTTPOnly:   utils.IntToBool(isHTTPOnly),
			HasExpire:    utils.IntToBool(hasExpire),
			IsPersistent: utils.IntToBool(isPersistent),
			CreateDate:   utils.TimeEpochFormat(createDate),
			ExpireDate:   utils.TimeEpochFormat(expireDate),
		}

		if osType == "mac" {
			value, _ = decryptChromeAES(secretKey, encryptValue)
		} else {
			value, _ = decryptWindowsChrome(secretKey, encryptValue)
		}
		cookie.Value = string(value)

		decryptedCookies = append(decryptedCookies, cookie)
	}

	return decryptedCookies
}

func decryptChromeAES(secretKey, encryptValue []byte) ([]byte, error) {
	if len(encryptValue) > 3 {
		if len(secretKey) == 0 {
			return nil, errSecurityKeyIsEmpty
		}
		var chromeIV = []byte{32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
		//var chromeIV = []byte{20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20}

		//version := string(encryptValue[:3])
		value, err := aes128CBCDecrypt(secretKey, chromeIV, encryptValue[3:])

		// if v20, remove first 32 bytes
		if true {
			value = value[32:]
		}

		return value, err
	} else {
		return nil, errDecryptFailed
	}
}

func decryptWindowsChrome(secretKey, encryptValue []byte) ([]byte, error) {
	if len(encryptValue) > 3 {
		block, _ := aes.NewCipher(secretKey)
		gcm, _ := cipher.NewGCM(block)

		nonce := encryptValue[3 : 3+12]
		version := string(encryptValue[:3])

		value, err := gcm.Open(nil, nonce, encryptValue[3+12:], nil)

		if version == "v20" {
			value = value[32:]
		}

		// if v20, remove first 32 bytes
		return value, err
	} else {
		return nil, errDecryptFailed
	}
}

func aes128CBCDecrypt(key, iv, encryptPass []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encryptLen := len(encryptPass)

	dst := make([]byte, encryptLen)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, encryptPass)
	dst = PKCS5UnPadding(dst)
	return dst, nil
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpad := int(src[length-1])
	return src[:(length - unpad)]
}
