package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"macCookies/pkg/types"
	"strconv"
)

func CookieQuickManager(decryptedCookies []types.Cookie) {
	var cookieQuickManagers []types.CookieQuickManager
	for _, Cookie := range decryptedCookies {
		cookieQuickManager := types.CookieQuickManager{
			PathRaw:           Cookie.Path,
			HostRaw:           "https://" + Cookie.Host + "/",
			ExpiresRaw:        strconv.FormatInt(Cookie.ExpireDate.Unix(), 10),
			ContentRaw:        Cookie.Value,
			NameRaw:           Cookie.KeyName,
			SameSiteRaw:       "no_restriction",
			ThisDomainOnlyRaw: "false",
			StoreRaw:          "firefox-default",
			FirstPartyDomain:  "",
			HTTPOnlyRaw:       "",
		}

		cookieQuickManagers = append(cookieQuickManagers, cookieQuickManager)
	}

	file, _ := json.MarshalIndent(cookieQuickManagers, "", "\t")
	_ = ioutil.WriteFile("macCookies.json", file, 0644)

	fmt.Println("\x1b[32m[+] successfully exported!\x1b[0m")
	fmt.Println("\x1b[32m[+] import macCookies.json to Firefox with CookieQuickManager\x1b[0m")
}
