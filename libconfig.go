package libconfig

import "os"

// lc is the default Parser for basic use.
// It uses "env" as the tag and `os.LookupEnv` for the lookup function.
var lc = Parser{
	Tag:      "env",
	LookupFn: os.LookupEnv,
}

// Get populates the config struct with values from the environment
func Get(config interface{}) error {
	return lc.Get(config)
}
