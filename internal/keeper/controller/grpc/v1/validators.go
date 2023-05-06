package v1

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

const (
	_missingField = "not set"

	DefaultMaxUsernameLength   = 128
	DefaultMaxSecretNameLength = 256

	DefaultMetadataLimit = 2 * 1024 * 1024

	DefaultDataLimit = 4 * 1024 * 1024
)

// validateUsername validates provided username.
func validateUsername(name string) (string, bool) {
	if name == "" {
		return _missingField, false
	}

	if len(name) > DefaultMaxUsernameLength {
		return fmt.Sprintf("should be <= %d characters", DefaultMaxUsernameLength), false
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

// validateSecretName validates provided secret name.
func validateSecretName(name string) (string, bool) {
	if name == "" {
		return _missingField, false
	}

	if len(name) > DefaultMaxSecretNameLength {
		return fmt.Sprintf("should be <= %d characters", DefaultMaxSecretNameLength), false
	}

	return "", true
}

// validateMetadata validates provided metadata.
func validateMetadata(metadata []byte) (string, bool) {
	if len(metadata) > DefaultMetadataLimit {
		return fmt.Sprintf("should be <= %d characters", DefaultMetadataLimit), false
	}

	return "", true
}

// validateSecretData validates provided secret data.
func validateSecretData(data []byte) (string, bool) {
	if len(data) == 0 {
		return _missingField, false
	}

	if len(data) > DefaultDataLimit {
		return fmt.Sprintf("should be <= %d characters", DefaultDataLimit), false
	}

	return "", true
}

// validateSecret validates provided secret and data.
func validateSecret(
	name string,
	metadata, data []byte,
) (*errdetails.BadRequest, bool) {
	br := &errdetails.BadRequest{}

	if reason, ok := validateSecretName(name); !ok {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "name",
			Description: reason,
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	if reason, ok := validateMetadata(metadata); !ok {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "metadata",
			Description: reason,
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	if reason, ok := validateSecretData(data); !ok {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       "data",
			Description: reason,
		}

		br.FieldViolations = append(br.FieldViolations, v)
	}

	if len(br.FieldViolations) == 0 {
		return nil, true
	}

	return br, false
}
