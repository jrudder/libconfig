package libconfig

import (
	"fmt"
	"reflect"
)

// ErrCannotParseEnv is returned if the variable cannot be parsed into the type
// expected by the struct field, e.g. parsing "500" into int8 will return this.
// This indicates that either the struct field is the wrong type or that the
// data in the env var is simply incompatible with the type.
type ErrCannotParseEnv struct {
	Because error
	Kind    reflect.Kind
	Key     string
	Value   string
}

// NewErrCannotParseEnv creates an ErrCannotParseEnv error
func NewErrCannotParseEnv(err error, k reflect.Kind, key, value string) *ErrCannotParseEnv {
	return &ErrCannotParseEnv{
		Because: err,
		Kind:    k,
		Key:     key,
		Value:   value,
	}
}

// Error returns a human-readable description of the error
func (e *ErrCannotParseEnv) Error() string {
	result := fmt.Sprintf("cannot parse env [%s] with value [%s] to kind [%s]", e.Key, e.Value, e.Kind)

	if e.Because != nil {
		result = fmt.Sprintf("%s: %s", result, e.Because.Error())
	}

	return result
}

// Cause returns the error that caused this ErrCannotParseEnv error
func (e *ErrCannotParseEnv) Cause() error {
	return e.Because
}

// ErrCannotSetKind is returned if the kind of a field cannot be set from a string.
// This indicates that the logic is missing from `setValueFromString` to handle
// the given reflect.Kind.
type ErrCannotSetKind struct {
	Kind reflect.Kind
}

// NewErrCannotSetKind creates a ErrCannotSetKind error
func NewErrCannotSetKind(k reflect.Kind) *ErrCannotSetKind {
	return &ErrCannotSetKind{
		Kind: k,
	}
}

// Error returns a human-readable description of the error
func (e *ErrCannotSetKind) Error() string {
	return fmt.Sprintf("cannot set kind [%s]", e.Kind.String())
}

// ErrDecodeFailure is returned by `Retrieve` if the value could not be decoded by the
// requested decoder
type ErrDecodeFailure struct {
	Key     string
	Value   string
	Type    string
	Because error
}

// NewErrDecodeFailure creates a ErrDecodeFailure error which wraps the error describing
// the cause of the failure
func NewErrDecodeFailure(err error, key, value, typ string) *ErrDecodeFailure {
	return &ErrDecodeFailure{
		Key:     key,
		Value:   value,
		Type:    typ,
		Because: err,
	}
}

// Error returns a human-readable description of the error
func (e *ErrDecodeFailure) Error() string {
	result := fmt.Sprintf("failed to decode var [%s] with value [%s] as [%s]", e.Key, e.Value, e.Type)

	if e.Because != nil {
		result = fmt.Sprintf("%s: %s", result, e.Because.Error())
	}

	return result
}

// Cause returns the error that caused the ErrDecodeFailure
func (e *ErrDecodeFailure) Cause() error {
	return e.Because
}

// ErrInvalidConfigType is returned if Get is called with a value that is not a pointer
// to a struct. It must be a pointer so that Get can modify the values. It must be a
// struct to have tagged fields.
type ErrInvalidConfigType struct {
	Type reflect.Type
}

// NewErrInvalidConfigType creates an ErrInvalidConfigType
func NewErrInvalidConfigType(t reflect.Type) *ErrInvalidConfigType {
	return &ErrInvalidConfigType{
		Type: t,
	}
}

// Error returns a human-readable description of the error
func (e *ErrInvalidConfigType) Error() string {
	return fmt.Sprintf("config must be pointer to struct but got %s", e.Type.String())
}

// ErrInvalidTagOption is returned if the struct field tag has an unsupported option.
type ErrInvalidTagOption struct {
	Tag       string
	BadOption string
}

// NewErrInvalidTagOption creates an ErrInvalidTagOption
func NewErrInvalidTagOption(tag, option string) *ErrInvalidTagOption {
	return &ErrInvalidTagOption{
		Tag:       tag,
		BadOption: option,
	}
}

// Error returns a human-readable description of the error
func (e *ErrInvalidTagOption) Error() string {
	return fmt.Sprintf("tag [%s] contains unsupported option [%s]", e.Tag, e.BadOption)
}

// ErrMissingNameTag is returned if the passed config struct field is tagged but no
// name is provided, e.g. `env:""`
type ErrMissingNameTag struct {
	Tag string
}

// NewErrMissingNameTag create an ErrMissingNameTag
func NewErrMissingNameTag(t string) *ErrMissingNameTag {
	return &ErrMissingNameTag{
		Tag: t,
	}
}

// Error ensures that ErrMissingNameTag conforms to the Error interface
func (e *ErrMissingNameTag) Error() string {
	return fmt.Sprintf("tagged field must be named but got [%s]", e.Tag)
}

// ErrOverflow is returned if a numeric reflect.Value cannot be set because it would result in an overflow
type ErrOverflow struct {
	Kind  reflect.Kind
	Key   string
	Value string
}

// NewErrOverflow creates an ErrOverflow
func NewErrOverflow(k reflect.Kind, key, value string) *ErrOverflow {
	return &ErrOverflow{
		Kind:  k,
		Key:   key,
		Value: value,
	}
}

// Error returns a human-readable description of the error
func (e *ErrOverflow) Error() string {
	return fmt.Sprintf("overflow detected trying to set field of kind [%s] to value [%s] for key [%s]", e.Kind.String(), e.Value, e.Key)
}

// ErrVarNotFound is returned if the given key is not found by the lookup function
type ErrVarNotFound struct {
	Key string
}

// NewErrVarNotFound creates a ErrVarNotFound error
func NewErrVarNotFound(key string) *ErrVarNotFound {
	return &ErrVarNotFound{
		Key: key,
	}
}

// Error returns a human-readable description of the error
func (e *ErrVarNotFound) Error() string {
	return fmt.Sprintf("var not found for key [%s]", e.Key)
}

// ErrNestedTags is returned if a tagged struct contains a tagged field, which, if supported, could
// result in unexpected behavior due to the parsing order of structs and struct fields
type ErrNestedTags struct {
	Field string
	Key   string
}

// NewErrNestedTags creates a ErrNestedTags error
func NewErrNestedTags(field, key string) *ErrNestedTags {
	return &ErrNestedTags{
		Field: field,
		Key:   key,
	}
}

// Error returns a human-readable description of the error
func (e *ErrNestedTags) Error() string {
	return fmt.Sprintf("field [%s] with key [%s] contains one or more nested subfields", e.Field, e.Key)
}
