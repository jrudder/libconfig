package libconfig_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/jrudder/libconfig"
)

func TestErrCannotParseEnv(t *testing.T) {
	cause := fmt.Errorf("some error")
	err := libconfig.NewErrCannotParseEnv(cause, reflect.Int, "key", "value")
	require.Equal(t, "cannot parse env [key] with value [value] to kind [int]: some error", err.Error(), "error string must match")
}

func TestErrCannotParseEnvWithoutCause(t *testing.T) {
	err := libconfig.NewErrCannotParseEnv(nil, reflect.Int, "key", "value")
	require.Equal(t, "cannot parse env [key] with value [value] to kind [int]", err.Error(), "error string must match")
}

func TestErrCannotParseEnvCause(t *testing.T) {
	expected := errors.New("some error")
	err := libconfig.NewErrCannotParseEnv(expected, reflect.Int, "key", "value")
	cause := errors.Cause(err)
	require.Equal(t, expected, cause, "ErrCannotParseEnv must have a cause")
}

func TestErrCannotSetKind(t *testing.T) {
	err := libconfig.NewErrCannotSetKind(reflect.Interface)
	require.Equal(t, "cannot set kind [interface]", err.Error(), "error string must match")
}

func TestErrDecodeFailure(t *testing.T) {
	cause := fmt.Errorf("some error")
	err := libconfig.NewErrDecodeFailure(cause, "key", "value", "base64")
	require.Equal(t, "failed to decode var [key] with value [value] as [base64]: some error", err.Error(), "error string must match")
}

func TestErrDecodeFailureWithoutCause(t *testing.T) {
	err := libconfig.NewErrDecodeFailure(nil, "key", "value", "base64")
	require.Equal(t, "failed to decode var [key] with value [value] as [base64]", err.Error(), "error string must match")
}

func TestErrDecodeFailureCause(t *testing.T) {
	expected := errors.New("some error")
	err := libconfig.NewErrDecodeFailure(expected, "key", "value", "base64")
	cause := errors.Cause(err)
	require.Equal(t, expected, cause, "ErrDecodeFailure must have a cause")
}

func TestErrInvalidConfigType(t *testing.T) {
	err := libconfig.NewErrInvalidConfigType(reflect.TypeOf(int(623)))
	require.Equal(t, "config must be pointer to struct but got int", err.Error(), "error string must match")
}

func TestErrInvalidConfigTypePtr(t *testing.T) {
	value := 623
	err := libconfig.NewErrInvalidConfigType(reflect.TypeOf(&value))
	require.Equal(t, "config must be pointer to struct but got *int", err.Error(), "error string must match")
}

func TestErrInvalidTagOption(t *testing.T) {
	err := libconfig.NewErrInvalidTagOption("tag,here", "something")
	require.Equal(t, "tag [tag,here] contains unsupported option [something]", err.Error(), "error string must match")
}

func TestErrMissingNameTag(t *testing.T) {
	err := libconfig.NewErrMissingNameTag("some-tag")
	require.Equal(t, "tagged field must be named but got [some-tag]", err.Error(), "error string must match")
}

func TestErrOverflow(t *testing.T) {
	err := libconfig.NewErrOverflow(reflect.Int8, "key", "value")
	require.Equal(t, "overflow detected trying to set field of kind [int8] to value [value] for key [key]", err.Error(), "error string must match")
}

func TestErrVarNotFound(t *testing.T) {
	err := libconfig.NewErrVarNotFound("key")
	require.Equal(t, "var not found for key [key]", err.Error(), "error string must match")
}

func TestErrNestedTags(t *testing.T) {
	err := libconfig.NewErrNestedTags("field", "key")
	require.Equal(t, "field [field] with key [key] contains one or more nested subfields", err.Error(), "error string must match")
}
