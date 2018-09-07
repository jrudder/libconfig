package libconfig_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/jrudder/libconfig"
)

func TestSingleton(t *testing.T) {
	key := "LIBCONFIG_SINGLETON"
	value := time.Now().String()
	os.Setenv(key, value)

	type Config struct {
		Value string `env:"LIBCONFIG_SINGLETON"`
	}

	config := Config{}
	err := libconfig.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(value, config.Value, "value should parse correctly")
}

func TestInvalidConfigTypeNotPointer(t *testing.T) {
	p := mapToParser(nil)

	var config int
	err := p.Get(config)
	expected := libconfig.NewErrInvalidConfigType(reflect.TypeOf(config))

	require := require.New(t)
	require.Equal(expected, err, "Get should fail with ErrInvalidConfigType")
}

func TestNoTags(t *testing.T) {
	type Config struct {
		VarA string
	}

	p := mapToParser(nil)

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
}

func TestTagMissingName(t *testing.T) {
	type Config struct {
		VarA string `env:""`
	}

	p := mapToParser(nil)

	config := Config{}
	err := p.Get(&config)
	expected := libconfig.NewErrMissingNameTag("")

	require := require.New(t)
	require.Equal(expected, err, "Get should fail")
}

func TestString(t *testing.T) {
	type Config struct {
		VarA string `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VAL_A",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal("VAL_A", config.VarA, "VarA should parse correctly")
}

func TestStringRequiredButMissing(t *testing.T) {
	type Config struct {
		VarA string `env:"VAR_A"`
	}

	p := mapToParser(nil)
	config := Config{}
	err := p.Get(&config)
	expected := libconfig.NewErrVarNotFound("VAR_A")

	require := require.New(t)
	require.Equal(expected, err, "Get should fail because VAR_A is not available")
}

func TestStringOptionalAndMissing(t *testing.T) {
	type Config struct {
		VarA string `env:"VAR_A"`
		VarB string `env:"VAR_B,optional"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VAL_A",
	})
	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get not should because VAR_B is marked as optional")
}

func TestBadOption(t *testing.T) {
	type Config struct {
		VarA string `env:"VAR_A,not-a-valid-option"`
	}

	p := mapToParser(nil)
	config := Config{}
	err := p.Get(&config)
	expected := libconfig.NewErrInvalidTagOption("VAR_A,not-a-valid-option", "not-a-valid-option")

	require := require.New(t)
	require.Equal(expected, err, "Get not should because VAR_B is marked as optional")
}

func TestByteSlice(t *testing.T) {
	type Config struct {
		VarA []byte `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VAL_A",
	})

	expected := []byte{'V', 'A', 'L', '_', 'A'}
	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(expected, config.VarA, "VarA should parse correctly")
}

func TestBase64String(t *testing.T) {
	type Config struct {
		VarA string `env:"VAR_A,base64"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VkFMX0E=",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal("VAL_A", config.VarA, "VarA should parse correctly")
}

func TestBase64ByteSlice(t *testing.T) {
	type Config struct {
		VarA []byte `env:"VAR_A,base64"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VkFMX0E=",
	})

	expected := []byte{'V', 'A', 'L', '_', 'A'}
	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(expected, config.VarA, "VarA should parse correctly")
}

func TestBase64Int(t *testing.T) {
	type Config struct {
		VarA int `env:"VAR_A,base64"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "MDkxNQ==",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(915, config.VarA, "VarA should parse correctly")
}

func TestBase64Invalid(t *testing.T) {
	type Config struct {
		VarA string `env:"VAR_A,base64"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "i-am-not-base64",
	})

	config := Config{}
	err := p.Get(&config)
	// Note that we do not actually expect a nil error.
	// We care (and test below) that an error is present, but not the error itself.
	expected := libconfig.NewErrDecodeFailure(nil, "VAR_A", "i-am-not-base64", "base64")

	require := require.New(t)
	require.Error(err, "Get should fail to parse the value as the base64")
	specificErr, ok := err.(*libconfig.ErrDecodeFailure)
	require.True(ok, "the error should be ErrDecodeFailure")
	require.Error(specificErr.Because, "Because should be set")
	specificErr.Because = nil // clear the underlying error so that we can validate the rest of the struct using `expected`
	require.Equal(expected, err, "Get should fail to parse the value as the kind")
}

func TestTwoStrings(t *testing.T) {
	type Config struct {
		VarA string `env:"VAR_A"`
		VarB string `env:"VAR_B"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VAL_A",
		"VAR_B": "VAL_B",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal("VAL_A", config.VarA, "VarA should parse correctly")
	require.Equal("VAL_B", config.VarB, "VarB should parse correctly")
}

func TestInt(t *testing.T) {
	type Config struct {
		VarA int `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "500",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(500, config.VarA, "VarA should parse correctly")
}

func TestIntPointer(t *testing.T) {
	type Config struct {
		VarA *int `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "500",
	})

	config := Config{}
	err := p.Get(&config)
	expected := 500
	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(&expected, config.VarA, "VarA should parse correctly")
}

func TestStringPointer(t *testing.T) {
	type Config struct {
		VarA *string `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "500",
	})

	config := Config{}
	err := p.Get(&config)
	expected := "500"
	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(&expected, config.VarA, "VarA should parse correctly")
}

func TestIntOverflow(t *testing.T) {
	type Config struct {
		VarA int8 `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "500",
	})

	config := Config{}
	err := p.Get(&config)

	expected := libconfig.NewErrOverflow(reflect.Int8, "VAR_A", "500")

	require := require.New(t)
	require.Equal(expected, err, "Get should fail to parse \"500\" as int8")
}

func TestIntCannotParseEnv(t *testing.T) {
	type Config struct {
		VarA int `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "not-an-int",
	})

	config := Config{}
	err := p.Get(&config)
	// Note that we do not actually expect a nil error.
	// We care (and test below) that an error is present, but not the error itself.
	expected := libconfig.NewErrCannotParseEnv(nil, reflect.Int, "VAR_A", "not-an-int")

	require := require.New(t)
	require.Error(err, "Get should fail to parse the value as the kind")
	specificErr, ok := err.(*libconfig.ErrCannotParseEnv)
	require.True(ok, "the error should be ErrCannotParseEnv")
	require.Error(specificErr.Because, "Because should be set")
	specificErr.Because = nil // clear the underlying error so that we can validate the rest of the struct using `expected`
	require.Equal(expected, err, "Get should fail to parse the value as the kind")
}

func TestUint(t *testing.T) {
	type Config struct {
		VarA uint `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "500",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(uint(500), config.VarA, "VarA should parse correctly")
}

func TestUintOverflow(t *testing.T) {
	type Config struct {
		VarA uint8 `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "500",
	})

	config := Config{}
	err := p.Get(&config)

	expected := libconfig.NewErrOverflow(reflect.Uint8, "VAR_A", "500")

	require := require.New(t)
	require.Equal(expected, err, "Get should fail to parse \"500\" as uint8")
}

func TestUintCannotParseEnv(t *testing.T) {
	type Config struct {
		VarA uint `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "not-a-uint",
	})

	config := Config{}
	err := p.Get(&config)
	// Note that we do not actually expect a nil error.
	// We care (and test below) that an error is present, but not the error itself.
	expected := libconfig.NewErrCannotParseEnv(nil, reflect.Uint, "VAR_A", "not-a-uint")

	require := require.New(t)
	require.Error(err, "Get should fail to parse the value as the kind")
	specificErr, ok := err.(*libconfig.ErrCannotParseEnv)
	require.True(ok, "the error should be ErrCannotParseEnv")
	require.Error(specificErr.Because, "Because should be set")
	specificErr.Because = nil // clear the underlying error so that we can validate the rest of the struct using `expected`
	require.Equal(expected, err, "Get should fail to parse the value as the kind")
}

func TestFloat32(t *testing.T) {
	type Config struct {
		VarA float32 `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "500.5",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(float32(500.5), config.VarA, "VarA should parse correctly")
}

func TestFloat64(t *testing.T) {
	type Config struct {
		VarA float64 `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "500.5",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(float64(500.5), config.VarA, "VarA should parse correctly")
}

func TestFloatOverflow(t *testing.T) {
	type Config struct {
		VarA float32 `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "1003006009012015018021024027030033036039",
	})

	config := Config{}
	err := p.Get(&config)

	expected := libconfig.NewErrOverflow(reflect.Float32, "VAR_A", "1003006009012015018021024027030033036039")

	require := require.New(t)
	require.Equal(expected, err, "Get should fail to parse \"1003006009012015018021024027030033036039\" as float32")
}

func TestFloatCannotParseEnv(t *testing.T) {
	type Config struct {
		VarA float64 `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "not-a-float",
	})

	config := Config{}
	err := p.Get(&config)
	// Note that we do not actually expect a nil error.
	// We care (and test below) that an error is present, but not the error itself.
	expected := libconfig.NewErrCannotParseEnv(nil, reflect.Float64, "VAR_A", "not-a-float")

	require := require.New(t)
	require.Error(err, "Get should fail to parse the value as the kind")
	specificErr, ok := err.(*libconfig.ErrCannotParseEnv)
	require.True(ok, "the error should be ErrCannotParseEnv")
	require.Error(specificErr.Because, "Because should be set")
	specificErr.Because = nil // clear the underlying error so that we can validate the rest of the struct using `expected`
	require.Equal(expected, err, "Get should fail to parse the value as the kind")
}

func TestBoolTrue(t *testing.T) {
	type Config struct {
		VarA bool `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "true",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(true, config.VarA, "VarA should parse correctly")
}

func TestBoolFalse(t *testing.T) {
	type Config struct {
		VarA bool `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "false",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal(false, config.VarA, "VarA should parse correctly")
}

func TestBoolCannotParseEnv(t *testing.T) {
	type Config struct {
		VarA bool `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "not-a-bool",
	})

	config := Config{}
	err := p.Get(&config)
	// Note that we do not actually expect a nil error.
	// We care (and test below) that an error is present, but not the error itself.
	expected := libconfig.NewErrCannotParseEnv(nil, reflect.Bool, "VAR_A", "not-a-bool")

	require := require.New(t)
	require.Error(err, "Get should fail to parse the value as the kind")
	specificErr, ok := err.(*libconfig.ErrCannotParseEnv)
	require.True(ok, "the error should be ErrCannotParseEnv")
	require.Error(specificErr.Because, "Because should be set")
	specificErr.Because = nil // clear the underlying error so that we can validate the rest of the struct using `expected`
	require.Equal(expected, err, "Get should fail to parse the value as the kind")
}
func TestErrCannotSetKindForInterface(t *testing.T) {
	type Config struct {
		VarA interface{} `env:"VAR_A"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "false",
	})

	config := Config{}
	err := p.Get(&config)
	expected := libconfig.NewErrCannotSetKind(reflect.Interface)

	require := require.New(t)
	require.Equal(expected, err, "Get should fail to parse reflect.Interface")
}
func TestStringsAndInts(t *testing.T) {
	type Config struct {
		VarA string `env:"VAR_A"`
		VarB string `env:"VAR_B"`
		VarC int    `env:"VAR_C"`
		VarD uint   `env:"VAR_D"`
		VarE int16  `env:"VAR_E"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VAL_A",
		"VAR_B": "VAL_B",
		"VAR_C": "10",
		"VAR_D": "20",
		"VAR_E": "30",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal("VAL_A", config.VarA, "VarA should parse correctly")
	require.Equal("VAL_B", config.VarB, "VarB should parse correctly")
	require.Equal(int(10), config.VarC, "VarC should parse correctly")
	require.Equal(uint(20), config.VarD, "VarD should parse correctly")
	require.Equal(int16(30), config.VarE, "VarE should parse correctly")
}

func TestNestedStruct(t *testing.T) {
	type Config struct {
		VarA   string `env:"VAR_A"`
		VarB   string `env:"VAR_B"`
		Nested struct {
			VarC int   `env:"VAR_C"`
			VarD uint  `env:"VAR_D"`
			VarE int16 `env:"VAR_E"`
		}
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VAL_A",
		"VAR_B": "VAL_B",
		"VAR_C": "10",
		"VAR_D": "20",
		"VAR_E": "30",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal("VAL_A", config.VarA, "VarA should parse correctly")
	require.Equal("VAL_B", config.VarB, "VarB should parse correctly")
	require.Equal(int(10), config.Nested.VarC, "VarC should parse correctly")
	require.Equal(uint(20), config.Nested.VarD, "VarD should parse correctly")
	require.Equal(int16(30), config.Nested.VarE, "VarE should parse correctly")
}

func TestNestedStructPointer(t *testing.T) {
	type Config struct {
		VarA   string `env:"VAR_A"`
		VarB   string `env:"VAR_B"`
		Nested *struct {
			VarC int   `env:"VAR_C"`
			VarD uint  `env:"VAR_D"`
			VarE int16 `env:"VAR_E"`
		}
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VAL_A",
		"VAR_B": "VAL_B",
		"VAR_C": "10",
		"VAR_D": "20",
		"VAR_E": "30",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal("VAL_A", config.VarA, "VarA should parse correctly")
	require.Equal("VAL_B", config.VarB, "VarB should parse correctly")
	require.Equal(int(10), config.Nested.VarC, "VarC should parse correctly")
	require.Equal(uint(20), config.Nested.VarD, "VarD should parse correctly")
	require.Equal(int16(30), config.Nested.VarE, "VarE should parse correctly")
}
func TestNestedStructError(t *testing.T) {
	type Config struct {
		VarA   string `env:"VAR_A"`
		VarB   string `env:"VAR_B"`
		Nested struct {
			VarC int `env:"VAR_C"`
		}
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VAL_A",
		"VAR_B": "VAL_B",
		"VAR_C": "not-an-int",
	})

	config := Config{}
	err := p.Get(&config)

	require := require.New(t)
	require.Error(err, "Get should fail")
}

func TestArrayAsJSON(t *testing.T) {
	type Config struct {
		VarA  string `env:"VAR_A"`
		VarB  string `env:"VAR_B"`
		Array []int  `env:"ARRAY,json"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VAL_A",
		"VAR_B": "VAL_B",
		"ARRAY": "[9,1,5]",
	})

	config := Config{}
	err := p.Get(&config)
	expected := []int{9, 1, 5}
	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal("VAL_A", config.VarA, "VarA should parse correctly")
	require.Equal("VAL_B", config.VarB, "VarB should parse correctly")
	require.Equal(expected, config.Array, "Array should parse correctly")
}

func TestArrayPointerAsJSON(t *testing.T) {
	type Config struct {
		VarA  string `env:"VAR_A"`
		VarB  string `env:"VAR_B"`
		Array []*int `env:"ARRAY,json"`
	}

	p := mapToParser(map[string]string{
		"VAR_A": "VAL_A",
		"VAR_B": "VAL_B",
		"ARRAY": "[9,1,5]",
	})

	config := Config{}
	err := p.Get(&config)
	a, b, c := 9, 1, 5
	expected := []*int{&a, &b, &c}
	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal("VAL_A", config.VarA, "VarA should parse correctly")
	require.Equal("VAL_B", config.VarB, "VarB should parse correctly")
	require.Equal(expected, config.Array, "Array should parse correctly")
}

func TestNestedStructPointerAsJSON(t *testing.T) {
	type Nested struct {
		VarC int    `json:"varc"`
		VarD string `json:"vard"`
	}
	type Config struct {
		VarA   string  `env:"VAR_A"`
		VarB   string  `env:"VAR_B"`
		Nested *Nested `env:"NESTED,json"`
	}

	p := mapToParser(map[string]string{
		"VAR_A":  "VAL_A",
		"VAR_B":  "VAL_B",
		"NESTED": `{"varc": 10, "vard": "val_d"}`,
	})

	config := Config{}
	err := p.Get(&config)
	expected := &Nested{VarC: 10, VarD: "val_d"}
	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal("VAL_A", config.VarA, "VarA should parse correctly")
	require.Equal("VAL_B", config.VarB, "VarB should parse correctly")
	require.Equal(expected, config.Nested, "Nested should parse correctly")
}

func TestNestedStructAsJSON(t *testing.T) {
	type Nested struct {
		VarC int    `json:"varc"`
		VarD string `json:"vard"`
	}
	type Config struct {
		VarA   string `env:"VAR_A"`
		VarB   string `env:"VAR_B"`
		Nested `env:"NESTED,json"`
	}

	p := mapToParser(map[string]string{
		"VAR_A":  "VAL_A",
		"VAR_B":  "VAL_B",
		"NESTED": `{"varc": 10, "vard": "val_d"}`,
	})

	config := Config{}
	err := p.Get(&config)
	expected := Nested{VarC: 10, VarD: "val_d"}
	require := require.New(t)
	require.NoError(err, "Get should not fail")
	require.Equal("VAL_A", config.VarA, "VarA should parse correctly")
	require.Equal("VAL_B", config.VarB, "VarB should parse correctly")
	require.Equal(expected, config.Nested, "Nested should parse correctly")
}
func TestNestedStructWithConfigTags(t *testing.T) {
	type Nested struct {
		VarC int `json:"varc" env:"VAR_C"`
	}
	type Config struct {
		Nested `env:"NESTED,json"`
	}

	p := mapToParser(map[string]string{
		"NESTED": "{}",
	})

	config := Config{}
	err := p.Get(&config)
	expected := libconfig.NewErrNestedTags("Nested", "NESTED")
	require := require.New(t)
	require.Equal(expected, err, "Get should fail because the struct is tagged and has tagged members")
}
func TestNestedStructAsInvalidJSON(t *testing.T) {
	type Nested struct {
		VarC int    `json:"varc"`
		VarD string `json:"vard"`
	}
	type Config struct {
		VarA   string `env:"VAR_A"`
		VarB   string `env:"VAR_B"`
		Nested `env:"NESTED,json"`
	}

	p := mapToParser(map[string]string{
		"VAR_A":  "VAL_A",
		"VAR_B":  "VAL_B",
		"NESTED": "i-am-not-json",
	})

	config := Config{}
	err := p.Get(&config)
	// Note that we do not actually expect a nil error.
	// We care (and test below) that an error is present, but not the error itself.
	expected := libconfig.NewErrDecodeFailure(nil, "NESTED", "i-am-not-json", "json")

	require := require.New(t)
	require.Error(err, "Get should fail to parse the value as JSON")
	specificErr, ok := err.(*libconfig.ErrDecodeFailure)
	require.True(ok, "the error should be ErrDecodeFailure")
	require.Error(specificErr.Because, "Because should be set")
	specificErr.Because = nil // clear the underlying error so that we can validate the rest of the struct using `expected`
	require.Equal(expected, err, "Get should fail to parse the value as the kind")
}

func mapToParser(envs map[string]string) libconfig.Parser {
	return libconfig.Parser{
		Tag: "env",
		LookupFn: func(name string) (string, bool) {
			value, found := envs[name]
			return value, found
		},
	}
}
