package libconfig

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
)

// Parser provides the core logic for libconfig.
// Typically, you will just use libconfig.Get, which uses a singleton
type Parser struct {
	Tag string

	// LookupFn enables the code to be thoroughly testable without relying on the
	// actual environment used during testing
	LookupFn func(key string) (string, bool)
}

// Get retrieves the configuration for the given struct by gathering values
// from the given LookupFn
func (p *Parser) Get(config interface{}) error {
	v := reflect.ValueOf(config)
	if t := v.Type(); !(t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct) {
		return NewErrInvalidConfigType(t)
	}

	_, err := p.parse(v.Elem())

	return err
}

// parse the given interface, looking for our tag, which indicates
// that the field can be populated by an environment variable
func (p *Parser) parse(config reflect.Value) (bool, error) {
	var tagFound bool

	// Look at each field of the struct
	t := config.Type()

	for i := 0; i < t.NumField(); i++ {
		// Get the struct field tag data
		field := t.Field(i)
		value := config.Field(i)
		tag, err := parseTag(field, p.Tag)
		if err != nil {
			return tagFound, err
		}

		// Parse tagged fields
		if tag.Tagged {
			tagFound = true

			// Get the value from the LookupFn
			err = p.retrieve(value, tag)
			if err != nil {
				return tagFound, err
			}
		}

		// If the field is a struct or pointer-to-struct, parse it
		if field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			// If the field is a pointer-to-struct, get the struct, not the pointer
			if field.Type.Kind() == reflect.Ptr {
				// If the pointer is nil, allocate memory first
				if value.IsNil() {
					value.Set(reflect.New(field.Type.Elem()))
				}
				value = value.Elem()
			}

			found, err := p.parse(value)

			// First ensure that a tagged struct contains no tagged members
			if tag.Tagged && found {
				return tagFound, NewErrNestedTags(field.Name, tag.Name)
			}

			// Handle any errors second
			if err != nil {
				return tagFound, err
			}
		}
	}

	return tagFound, nil
}

// retrieve gets the value for the tag from the lookup function, handling any
// necessary decoding, such as base64.
func (p *Parser) retrieve(v reflect.Value, tag tagData) error {
	var bytes []byte
	var err error

	value, found := p.LookupFn(tag.Name)
	if !found {
		if !tag.Optional {
			return NewErrVarNotFound(tag.Name)
		}

		return nil
	}

	// Base64-decode if specified
	if tag.Base64 {
		bytes, err = base64.StdEncoding.DecodeString(value)
		if err != nil {
			return NewErrDecodeFailure(err, tag.Name, value, "base64")
		}
	} else {
		bytes = []byte(value)
	}

	// JSON-decode if specified
	if tag.JSON {
		if v.Kind() == reflect.Ptr {
			// If v is a nil pointer, we need to allocate memory
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
		} else {
			// We need a pointer to the struct for unmarshalling
			v = v.Addr()
		}

		err = json.Unmarshal(bytes, v.Interface())
		if err != nil {
			return NewErrDecodeFailure(err, tag.Name, value, "json")
		}

		return nil
	}

	if v.Kind() == reflect.Ptr {
		// v is a Pointer; we need to allocate memory
		v.Set(reflect.New(v.Type().Elem()))
		v = v.Elem()
	}

	err = setValue(v, tag.Name, bytes)

	return err
}
