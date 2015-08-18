package envflag

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	debugEnabled = false
)

// EnvFlag represents a envflag object that contains several settings.
type EnvFlag struct {
	flagSet           *flag.FlagSet     // the flagSet to process
	minLength         int               // minLength defines min length of flag key, in order to support shortcut
	envFlagDict       map[string]string // envFlagDict is a user-defined env-flag map
	showEnvKeyInUsage bool              // show env key in usage
	showEnvValInUsage bool              // show env value in usage as default flag value
}

// ProcessFlagWithEnv parses flag from env.
// This function also add ENVIRONMENT VARIABLE to usage.
// NOTICE: flag.Parse() will not be called by this function.
func (ef EnvFlag) ProcessFlagWithEnv() error {
	if ef.flagSet.Parsed() {
		return errors.New("flag has already been parsed")
	}
	flagEnvMap := map[string]string{}
	for k, v := range ef.envFlagDict {
		flagEnvMap[v] = k
	}
	// Rewrite flag.Useage and show env key in usage
	if ef.showEnvKeyInUsage {
		ef.flagSet.VisitAll(func(f *flag.Flag) {
			if len(f.Name) < ef.minLength {
				return
			}
			envKey, ok := flagEnvMap[f.Name]
			if !ok {
				envKey = flagToEnv(f.Name)
			}
			envPrefix := fmt.Sprintf("[%s]", envKey)
			if strings.HasPrefix(f.Usage, envPrefix) {
				return
			}
			f.Usage = fmt.Sprintf("%s %s", envPrefix, f.Usage)
		})
	}

	for _, envLine := range os.Environ() {
		debug("Find a new line of environment variable, ", envLine)
		envKV := strings.SplitN(envLine, "=", 2)
		var key, value string
		key = envKV[0]
		if len(key) < ef.minLength {
			continue
		}
		if len(envKV) > 1 {
			value = envKV[1]
		} else {
			value = ""
		}

		var flagKey string
		if userFlag, ok := ef.envFlagDict[key]; ok {
			flagKey = userFlag
		} else {
			flagKey = envToFlag(key)
		}
		debug("ENV ", key, " is converted to ", flagKey)
		f := ef.flagSet.Lookup(flagKey)
		if f == nil {
			debug(flagKey, " is not defined in flag, skip!")
			continue
		}

		if ef.showEnvValInUsage {
			// use env value as default value (in usage)
			f.DefValue = value
		}

		debug("Set  [", flagKey, ",", value, "] to flag")
		if err := ef.flagSet.Set(flagKey, value); err != nil {
			return fmt.Errorf("error when set [%s,%s] into flag\n", flagKey, value)
		}
	}
	return nil
}

// ProcessFlagWithEnv parses flag from env.
// This function also add ENVIRONMENT VARIABLE to usage.
// NOTICE: flag.Parse() will not be called by this function.
func ProcessFlagWithEnv() error {
	return std.ProcessFlagWithEnv()
}

// Parse parses flag definitions from env and the argument list.
// Value from env can be overrided by the argument list.
// This function also add ENVIRONMENT VARIABLE to usage.
// It is extremely recommended to call this function in main()
// after all flags are defined and before flags are accessed by the program.
// NOTICE: flag.Parse() will be called by this function.
func (ef EnvFlag) Parse(arguments []string) error {
	if err := ef.ProcessFlagWithEnv(); err != nil {
		return err
	}
	return ef.flagSet.Parse(arguments)
}

// Parse parses the command-line flags from env and os.Args[1:].
// Value from env can be overrided by os.Args[1:].
// This function also add ENVIRONMENT VARIABLE to usage.
// It is extremely recommended to call this function in main()
// after all flags are defined and before flags are accessed by the program.
// NOTICE: flag.Parse() will be called by this function.
func Parse() error {
	return std.Parse(os.Args[1:])
}

// SetMinLength sets the min length.
// EnvFlag only parses the environment variables that is longer than min length
// and modify usage that is longer than min length.
func (ef *EnvFlag) SetMinLength(v int) {
	ef.minLength = v
}

// SetMinLength sets the min length for standard envflag.
// EnvFlag only parses the environment variables that is longer than min length
// and modify usage that is longer than min length.
func SetMinLength(v int) {
	std.SetMinLength(v)
}

// SetEnvFlagDict sets a user-defined env-flag map.
func (ef *EnvFlag) SetEnvFlagDict(v map[string]string) {
	ef.envFlagDict = v
}

// SetEnvFlagDict sets a user-defined env-flag map for standard envflag.
func SetEnvFlagDict(v map[string]string) {
	std.SetEnvFlagDict(v)
}

// SetShowEnvKeyInUsage sets whether show env key in usage.
func (ef *EnvFlag) SetShowEnvKeyInUsage(v bool) {
	ef.showEnvKeyInUsage = v
}

// SetShowEnvKeyInUsage sets whether show env key in usage for standard envflag.
func SetShowEnvKeyInUsage(v bool) {
	std.showEnvKeyInUsage = v
}

// SetShowEnvValInUsage sets whether show env value in usage.
// This is useful for confirming the environment variables.
func (ef *EnvFlag) SetShowEnvValInUsage(v bool) {
	ef.showEnvValInUsage = v
}

// SetShowEnvValInUsage sets whether show env value in usage for standard envflag.
// This is useful for confirming the environment variables.
func SetShowEnvValInUsage(v bool) {
	std.showEnvValInUsage = v
}

var std = NewEnvFlag(flag.CommandLine, 3, map[string]string{}, true, true)

// NewEnvFlag returns a new EnvFlag.
func NewEnvFlag(
	flagSet *flag.FlagSet,
	minLength int,
	envFlagDict map[string]string,
	showEnvKeyInUsage bool,
	showEnvValInUsage bool) *EnvFlag {
	return &EnvFlag{
		flagSet:           flagSet,
		minLength:         minLength,
		envFlagDict:       envFlagDict,
		showEnvKeyInUsage: showEnvKeyInUsage,
		showEnvValInUsage: showEnvValInUsage,
	}
}

// DebugEnabled returns whether the debug is enabled or not
func DebugEnabled() bool {
	return debugEnabled
}

// SetDebugEnabled enables debug info
func SetDebugEnabled(enabled bool) {
	debugEnabled = enabled
}

// envToFlag converts THIS_FORMAT to this-format
func envToFlag(e string) string {
	return strings.Replace(strings.ToLower(e), "_", "-", -1)
}

// flagToEnv converts this-format to THIS_FORMAT
func flagToEnv(f string) string {
	return strings.Replace(strings.ToUpper(f), "-", "_", -1)
}

func debug(v ...interface{}) {
	if !debugEnabled {
		return
	}
	log.Println(v...)
}
