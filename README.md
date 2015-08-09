# envflag

[![License: APACHE2](http://img.shields.io/badge/license-APACHE2-yellow.svg)](LICENSE)

A simple golang tools to set flag via environment variables inspired by [Go: Best Practices for Production Environments](http://peter.bourgon.org/go-in-production/#configuration)

# Features

- Set flag via environment variables.
- Auto mapping environment variables to flag. (e.g. `DATABASE_PORT` to `-database-port`)
- Customizable env - flag mapping support.
- Automatically add environment variables to help `-h`.
- Min length (default is 3) support in order to avoid parsing short flag.

# Basic Usage

### __Just keep it SIMPLE and SIMPLE and SIMPLE!__

Use `envflag.Parse()` instead of `flag.Parse()`.

Here is an example.

```go
package main

import (
    "flag"
    "fmt"

    "github.com/kouhin/envflag"
)

func main() {
    var (
        databaseMasterHost = flag.String("database-master-host", "localhost", "Database master host")
        databaseMasterPort = flag.Int("database-master-username", 3306, "Database master port")
    )
    envflag.Parse()
    fmt.Println("RESULT: ", *databaseMasterHost, ":", *databaseMasterPort)
}
```

Run `go run main.go -h` you will get the following usage:

```
Usage of XXXX
  -database-master-host="localhost": [DATABASE_MASTER_HOST] Database master host
  -database-master-username=3306: [DATABASE_MASTER_USERNAME] Database master port
```

```
$go run main.go
RESULT:  localhost : 3306

$export DATABASE_MASTER_HOST=192.168.0.2
$go run main.go
RESULT:  192.168.0.2 : 3306
$go run main.go -database-master-host=192.168.0.3
RESULT:  192.168.0.3 : 3306
```

# Advanced Usage

You can customize envflag by `env.Config`[Optional].

```go
func main() {
    envflag.Setup(&envflag.Config{
        DebugEnabled: true, // Debug
        MinLength:    5, // Min length of environment variables
        EnvFlagDict: map[string]string{ // User-defined env-flag map
            "MY_APP_ENV": "app-env",
        },
    })
    var (
        appEnv = flag.String("app-env", "dev", "Application env")
    )
    envflag.Parse()
    fmt.Println("appEnv:", appEnv)
}
```
