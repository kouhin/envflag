package envflag

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Config provides configuration for this package
type Config struct {
	// DebugEnabled enables debug mode
	DebugEnabled bool
	// MinLength defines min length of flag key, in order to support shortcut
	MinLength int
	// EnvFlagDict is a user-defined env-flag map
	EnvFlagDict map[string]string
}

// DefaultConfig is the default configuration of envflag
var DefaultConfig = &Config{
	DebugEnabled: false,
	MinLength:    3,
	EnvFlagDict:  map[string]string{},
}

// config is the configuration of envflag package
var config = &Config{}

// Setup sets up this package with customized config
func Setup(c *Config) {
	if c.DebugEnabled {
		config.DebugEnabled = c.DebugEnabled
	} else {
		config.DebugEnabled = DefaultConfig.DebugEnabled
	}

	if c.MinLength != 0 {
		config.MinLength = c.MinLength
	} else {
		config.MinLength = DefaultConfig.MinLength
	}

	if c.EnvFlagDict != nil {
		config.EnvFlagDict = c.EnvFlagDict
	} else {
		config.EnvFlagDict = DefaultConfig.EnvFlagDict
	}
}

// envToFlag converts THIS_FORMAT to this-format
func envToFlag(e string) string {
	return strings.Replace(strings.ToLower(e), "_", "-", -1)
}

// flagToEnv converts this-format to THIS_FORMAT
func flagToEnv(f string) string {
	return strings.Replace(strings.ToUpper(f), "-", "_", -1)
}

// Parse parses the command-line flags from env and os.Args[1:].
// Value from env can be overrided by os.Args[1:].
// This function also add ENVIRONMENT VARIABLE to usage.
// It is extremely recommended to call this function in main()
// after all flags are defined and before flags are accessed by the program.
func Parse() {
	if flag.Parsed() {
		return
	}

	flagEnvMap := map[string]string{}
	for k, v := range config.EnvFlagDict {
		flagEnvMap[v] = k
	}
	// Rewrite flag.Useage
	flag.VisitAll(func(f *flag.Flag) {
		if len(f.Name) < config.MinLength {
			return
		}
		envKey, ok := flagEnvMap[f.Name]
		if !ok {
			envKey = flagToEnv(f.Name)
		}
		f.Usage = fmt.Sprintf("[%s] %s", envKey, f.Usage)
	})

	for _, envLine := range os.Environ() {
		debug("Find a new line of environment variable, ", envLine)
		envKV := strings.SplitN(envLine, "=", 2)
		var key, value string
		key = envKV[0]
		if len(key) < config.MinLength {
			continue
		}
		if len(envKV) > 1 {
			value = envKV[1]
		} else {
			value = ""
		}

		var flagKey string
		if userFlag, ok := config.EnvFlagDict[key]; ok {
			flagKey = userFlag
		} else {
			flagKey = envToFlag(key)
		}
		debug("ENV ", key, " is converted to ", flagKey)
		if flag.Lookup(flagKey) == nil {
			debug(flagKey, " is not defined in flag, skip!")
			continue
		}

		debug("Set  [", flagKey, ",", value, "] to flag")
		if err := flag.Set(flagKey, value); err != nil {
			log.Printf("error when set [%s,%s] into flag\n", flagKey, value)
		}
	}
	flag.Parse()
}

func debug(v ...interface{}) {
	if !config.DebugEnabled {
		return
	}
	log.Println(v)
}
