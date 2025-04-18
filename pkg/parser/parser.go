package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"time"

	"github.com/kawakatz/macCookies/pkg/types"
)

func CookieQuickManager(decryptedCookies []types.Cookie) {
	var cookieQuickManagers []types.CookieQuickManager
	for _, Cookie := range decryptedCookies {
		cookieQuickManager := types.CookieQuickManager{
			PathRaw:           Cookie.Path,
			HostRaw:           "https://" + Cookie.Host + "/",
			ExpiresRaw:        strconv.FormatInt(time.Now().AddDate(1, 0, 0).Unix(), 10),
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

func StorageAce(decryptedCookies []types.Cookie) {
	var storageAces []types.StorageAce
	for _, Cookie := range decryptedCookies {
		t := time.Now().AddDate(1, 0, 0)
		ts := float64(t.UnixNano()) / 1e9
		ts = math.Round(ts*1e6) / 1e6

		storageAce := types.StorageAce{
			Domain:         Cookie.Host,
			ExpirationDate: ts,
			HostOnly:       false,
			HTTPOnly:       false,
			Name:           Cookie.KeyName,
			Path:           Cookie.Path,
			SameSite:       "no_restriction",
			Secure:         false,
			Session:        false,
			StoreID:        "0",
			Value:          Cookie.Value,
		}

		storageAces = append(storageAces, storageAce)
	}

	file, _ := json.MarshalIndent(storageAces, "", "\t")
	_ = ioutil.WriteFile("macCookies.json", file, 0644)

	fmt.Println("\x1b[32m[+] successfully exported!\x1b[0m")
	fmt.Println("\x1b[32m[+] import macCookies.json to Google Chrome with StorageAce\x1b[0m")
}
