package main

import (
	"flag"
	"fmt"

	"github.com/kouhin/envflag"
)

func main() {
	var (
		databaseMasterHost = flag.String("database-master-host", "", "Database master host")
		databaseMasterPort = flag.Int("database-master-port", -1, "Database master port")
		mh                 = flag.String("mh", "", "a shortcut for database-master-host")
	)
	envflag.Parse()
	fmt.Println("RESULT: ", *databaseMasterHost, ":", *databaseMasterPort, mh)
}
