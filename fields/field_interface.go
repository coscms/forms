package fields

import (
	"github.com/coscms/forms/config"
)

// FieldInterface defines the interface an object must implement to be used in a form. Every method returns a FieldInterface object
// to allow methods chaining.
type FieldInterface = config.FieldInterface
