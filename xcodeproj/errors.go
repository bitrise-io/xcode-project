package xcodeproj

import "fmt"

// ConfigurationNotFoundError ...
type ConfigurationNotFoundError struct {
	name string
}

// Error ...
func (e ConfigurationNotFoundError) Error() string {
	return fmt.Sprintf("configuration not found with name: %s", e.name)
}

// NewConfigurationNotFoundError ...
func NewConfigurationNotFoundError(name string) ConfigurationNotFoundError {
	return ConfigurationNotFoundError{name: name}
}
