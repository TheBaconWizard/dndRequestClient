# dndRequestClient
A go tls client based off of [Carcraftz](https://github.com/Carcraftz)'s [TLS API](https://github.com/Carcraftz/TLS-Fingerprint-API)

# Example
```go
package main

import (
	"fmt"
	https "github.com/TheBaconWizard/dndRequestClient"
)

func main() {
	requestResponse, responseString := https.GetProxyless("https://google.com", map[string]string{})
	fmt.Println(responseString)
}
```

# Notes
- Default methods include Get, Post, and Patch. If you wish to use other methods you can use HandleReq.
- If you would like to easily format your headers you can use [this tool](https://www.connorstevens.dev/headers) made by [@cnrstvns](https://twitter.com/cnrstvns)

# Changelog
- [2/14/22] - Added proxy support if you want to continue without proxies change your methods to Proxyless, updated some documentation, fixed a bug where you would get nil pointers when making lots of requests at once.
- [12/29/21] - If you don't wish to use tls (maybe you don't want the handshake for simplicity) change your input url to "http://" this will enable nontls mode however your url will be re prepended to "https://" due to some webservers refusing non https traffic regardless of whether tls handshake is needed
# Contributing
If you have anything you wish to contribute feel free to open an issue or pr
