package libconfig

import (
	"reflect"
	"strings"
)

type tagData struct {
	Tagged   bool
	Name     string
	Optional bool
	Base64   bool
	JSON     bool
}

func parseTag(f reflect.StructField, tag string) (tagData, error) {
	result := tagData{}

	// Get the tags
	var tags string
	tags, result.Tagged = f.Tag.Lookup(tag)
	if !result.Tagged {
		return result, nil
	}

	// Split into tokens and then parse the tokens
	tagTokens := strings.Split(tags, ",")

	// Parse: Name
	result.Name = tagTokens[0]
	if len(result.Name) == 0 {
		return result, NewErrMissingNameTag(tags)
	}

	for i := 1; i < len(tagTokens); i++ {
		switch tagTokens[i] {
		case "base64":
			result.Base64 = true
		case "json":
			result.JSON = true
		case "optional":
			result.Optional = true
		default:
			return tagData{}, NewErrInvalidTagOption(tags, tagTokens[i])
		}
	}

	return result, nil
}
