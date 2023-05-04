package v1

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

const (
	_missingField             = "not set"
	_defaultMaxUsernameLength = 128
)

// validateUsername validates provided username.
func validateUsername(username string) (string, bool) {
	if username == "" {
		return _missingField, false
	}

	if len(username) > _defaultMaxUsernameLength {
		return fmt.Sprintf("should be <= %d characters", _defaultMaxUsernameLength), false
	}

	return "", true
}

// validateSecurityKey validates provided security key.
func validateSecurityKey(key string) (string, bool) {
	if key == "" {
		return _missingField, false
	}

	return "", true
}

// validateCredentials validates provided credentials.
func validateCredentials(username, key string) (*errdetails.BadRequest, bool) {
	br := &errdetails.BadRequest{}

	if reason, ok := validateUsername(username); !ok {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "username",
			Description: reason,
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	if reason, ok := validateSecurityKey(key); !ok {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "security_key",
			Description: reason,
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	if len(br.FieldViolations) == 0 {
		return nil, true
	}

	return br, false
}
