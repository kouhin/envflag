# envflag

[![License: APACHE2](https://img.shields.io/github/license/kouhin/envflag.svg)](LICENSE)

A simple golang tool to set flag via environment variables inspired by [Go: Best Practices for Production Environments](http://peter.bourgon.org/go-in-production/#configuration)

# Features

- Set flag via environment variables.
- Auto mapping environment variables to flag. (e.g. `DATABASE_PORT` to `-database-port`)
- Customizable env - flag mapping support.
- Min length (default is 3) support in order to avoid parsing short flag.
- Show environment variable key in usage (-h).
- Show environment variable value in usage as default value (-h) in order to confirm enviroment settings.

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
        databaseMasterPort = flag.Int("database-master-port", 3306, "Database master port")
    )
    if err := envflag.Parse(); err != nil {
	    panic(err)
	}
    fmt.Println("RESULT: ", *databaseMasterHost, ":", *databaseMasterPort)
}
```

Run `DATABASE_MASTER_HOST=192.168.0.2 go run main.go -h` you will get the following usage:

```
Usage of XXXX
  -database-master-host="192.168.0.2": [DATABASE_MASTER_HOST] Database master host
  -database-master-port=3306: [DATABASE_MASTER_PORT] Database master port
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

You can customize envflag [Optional].

```go
func main() {
    ef := envflag.NewEnvFlag(
	    flag.CommandLine, // which FlagSet to parse
		2, // min length
		map[string]string{ // User-defined env-flag map
            "MY_APP_ENV": "app-env",
        },
		true, // show env variable key in usage
		true, // show env variable value in usage
    )
    var (
        appEnv = flag.String("app-env", "dev", "Application env")
    )
    if err := ef.Parse(os.Args[1:]); err != nil {
	    panic(err)
	}
    fmt.Println("appEnv:", appEnv)
}
```

# Enable debug info

Use `envflag.SetDebugEnabled(true)`.
