package main

import (
	"encoding/base64"
	"flag"
	"strings"

	"github.com/kawakatz/macCookies/pkg/decrypt"
	"github.com/kawakatz/macCookies/pkg/parser"
	"github.com/kawakatz/macCookies/pkg/types"
)

func main() {
	osPtr := flag.Bool("win", false, "decrypt Windows' Cookies")
	stracePtr := flag.Bool("storageace", false, "parse Cookies for StoraeAce")
	flag.Parse()

	cmdArgs := flag.Args()
	browserName := cmdArgs[0]
	cookiesFile := cmdArgs[1]

	var decryptedCookies []types.Cookie
	if strings.EqualFold(browserName, "Firefox") {
		decryptedCookies = decrypt.FirefoxCookies(cookiesFile)
	} else if strings.EqualFold(browserName, "Safari") {
		decryptedCookies = decrypt.SafariCookies(cookiesFile)
	} else {
		var secretKey []byte
		if !*osPtr {
			//fmt.Println("mac")
			browserPassword := cmdArgs[2]
			secretKey = decrypt.MacPassword2SecretKey(browserPassword)
			decryptedCookies = decrypt.ChromeCookies(cookiesFile, secretKey, "mac")
		} else {
			//fmt.Println("win")
			secretKeyBase64 := cmdArgs[2]
			secretKey, _ = base64.StdEncoding.DecodeString(secretKeyBase64)
			decryptedCookies = decrypt.ChromeCookies(cookiesFile, secretKey, "win")
		}
	}

	if *stracePtr {
		parser.StorageAce(decryptedCookies)
	} else {
		parser.CookieQuickManager(decryptedCookies)
	}
}
