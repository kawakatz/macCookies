# macCookies🍪
<p align="center">
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-_red.svg"></a>
<a href="https://github.com/kawakatz/grepfiles/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>
<a href="https://twitter.com/kawakatz"><img src="https://img.shields.io/twitter/follow/kawakatz.svg?logo=twitter"></a>
</p>

macCookies decrypt cookies stored in macOS browsers for pentesters.<br>
This tool is intended to be used with C2.

# Installation
```sh
# install grepfiles
➜  ~ go install -v github.com/kawakatz/macCookies/cmd/macCookies@latest
```

## Usage
### Safari
- FDA (including Finder automation permission) is required to access Cookies.binarycookies
- Cookies.binarycookies is not encrypted

```sh
➜  ~ macCookies Safari ~/Library/Containers/com.apple.Safari/Data/Library/Cookies/Cookies.binarycookies
```

### Firefox
- cookies.sqlite is not encrypted

```sh
➜  ~ macCookies Firefox ~/Library/Application\ Support/Firefox/Profiles/<profile>/cookies.sqlite
```

### Google Chrome, Microsoft Edge, Slack Application, etc...
- login-keychain password is required to decrypt login-keychain

```sh
# extract Chrome Safe Storage value
➜  ~ ./chainbreaker.py --dump-all login.keychain-db --password=<login-keychain password>
➜  ~ macCookies Chrome ~/Library/Application\ Support/Google/Chrome/Default/Cookies <Chrome Safe Storage>
```

#### Notes
If the victim had downloaded the app from the AppStore, files that store Cookies is located under `~/Library/Containers/<bundle id>/Data/Library/Application Support/` because the app must be sandboxed.

If you do not know the password for login-keychain, you can use <a href="https://github.com/kawakatz/macCookieStealer">macCookieStealer</a> to retrieve cookies from chromium-based browsers.

There are also cases where it is possible to bypass keychain client validation by injecting the Dynamic Library into an older application, thereby taking the encryption key from the keychain. Since Google Chrome has long been built with the restrict flag, Dynamic Library injection is not possible and this technique is not effective.

#### Option
It is also possible to decrypt Cookies retrieved from Windows.
In that case, use <a href="https://github.com/crypt0p3g/bof-collection/tree/main/ChromiumKeyDump">ChromiumKeyDump</a> to retrieve the encryption key.
```sh
➜  ~ macCookies -win Chrome Cookies <encryption key>
```

## References
- https://github.com/cixtor/binarycookies (MIT License)<br>
    decryption logic for Safari
- https://github.com/moonD4rk/HackBrowserData (MIT License)<br>
    decryption logic for FIrefox, Google Chrome, Microsoft Edge, etc...