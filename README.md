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
	requestResponse, responseString := https.Get("https://google.com", map[string]string{})
	fmt.Println(responseString)
}
```

# Notes
- Default methods include Get, Post, and Patch. If you wish to use other methods you can use HandleReq.
- If you would like to easily format your headers you can use [this tool](https://www.connorstevens.dev/headers) made by [@cnrstvns](https://twitter.com/cnrstvns)

# Contributing
If you have anything you wish to contribute feel free to open an issue or pr
