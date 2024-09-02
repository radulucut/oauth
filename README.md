# oauth

Go (golang) client library for validating access tokens and getting basic user data (id, email, name) from a provider.

Providers: Google, Facebook, Microsoft.

## Install

`go get github.com/radulucut/oauth`

## Usage

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/radulucut/oauth/v2"
)

func main() {
    client := oauth.NewClient(oauth.Client{
        Timeout: 10 * time.Second,
    })
    res, err := client.Google("<access_token>")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%+v\n", res)
}
```
