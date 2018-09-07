// Package libconfig provides a simple method for populating a struct with data
// from environment variables. Struct tags specify the name of the environment variable
// and whether the value is base64 encoded, json, or optional (by default, anything
// tagged is required and an error will be returned if the corresponding environment
// variable is not set.
//
// The following basic example is enough to get started:
//
//   package main
//
//   // config declares our configuration struct and maps
//   // individual fields to environment variables
//   type config struct {
//       ConnectionString string `env:"CONN_STRING"`
//       LogLevel         int    `env:"LOG_LEVEL"`
//   }
//
//   func main() {
//       // Create our config struct, supplying any defaults
//       c := config{
//           LogLevel: 2,
//       }
//
//       // Get populates our config struct from environment variables
//       _ = libconfig.Get(os.LookupEnv, &c)
//
//       fmt.Printf("DB_URL: %s\n", c.ConnectionString)
//   }
//
// The field tag must begin with the environment variable name and may be followed
// by zero or more of: base64, json, and optional.
//
//   type Config struct {
//       // Basic parsing just need a name
//       BasicString string `env:"BASIC_STRING"`
//
//       // Parsing of basic types uses strconv.Parse*
//       BasicInt int `env:"BASIC_INT"`
//
//       // Since it is marked as optional, IntPtr will be nil if INT_PTR is unset
//       IntPtr *int `env:"INT_PTR,optional"`
//
//       // Values can be base64-encoded. Tagging with "base64" will cause libconfig to
//       // decode the string value prior to further parsing, so you can have a base64-encoded
//       // string, []byte, float32, etc.
//       StringFromB64 string `env:"BASE64_STRING,base64"`
//
//       // Anything that can be parsed can be base64-encoded
//       Float32FromB64 string `env:"BASE64_FLOAT32,base64"`
//
//       // Use JSON for structs
//       FromJSONStruct struct {
//           NestedOne string `json:"nested_one"`
//           NestedTwo uint32 `json:"nested_two"`
//       } `env:"JSON_STRUCT_DATA,json"`
//
//       // Use JSON for slices (except []byte, which can be parsed directly)
//       FromJSONArray []int `env:"JSON_INT_ARRAY,json"`
//
//       // Base64 and JSON can be used together
//       FromB64JSON string `env:"B64_JSON,base64,json"`
//
//       // Tag ordering does not matter, base64-decoding happens first if-specified
//       // (since that's the only reasonable option)
//       FromB64JSONAlso string `env:"B64_JSON,json,base64"`
//   }
//
// To use a different tag name, instead of the default of "env", create a Parser.
//
//   p := libconfig.Parser{
//       Tag: "envtag",
//       LookupFn: os.LookupEnv,
//   }
//
//   err := p.Get(&config)
//
package libconfig
