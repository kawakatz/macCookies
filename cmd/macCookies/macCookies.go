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
	flag.Parse()

	cmdArgs := flag.Args()
	browserName := cmdArgs[0]
	cookiesFile := cmdArgs[1]

	if strings.EqualFold(browserName, "Firefox") {
		decryptedCookies := decrypt.FirefoxCookies(cookiesFile)
		parser.CookieQuickManager(decryptedCookies)
	} else if strings.EqualFold(browserName, "Safari") {
		decryptedCookies := decrypt.SafariCookies(cookiesFile)
		parser.CookieQuickManager(decryptedCookies)
	} else {
		var secretKey []byte
		var decryptedCookies []types.Cookie
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

		parser.CookieQuickManager(decryptedCookies)
	}
}
