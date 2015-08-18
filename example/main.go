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
		mh                 = flag.String("mh", "localhost", "a shortcut for database-master-host")
	)
	envflag.SetDebugEnabled(true)
	if err := envflag.Parse(); err != nil {
		panic(err)
	}
	fmt.Println("RESULT: ", *databaseMasterHost, ":", *databaseMasterPort, *mh)
}
