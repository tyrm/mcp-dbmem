package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlreadyExistsError(t *testing.T) {
	t.Parallel()

	msg := "entry already exists"
	err := NewErrAlreadyExists(msg)
	assert.IsType(t, &AlreadyExistsError{}, err)
	assert.EqualError(t, err, msg)

	// Test Error() method directly
	customErr := &AlreadyExistsError{message: "custom error"}
	assert.Equal(t, "custom error", customErr.Error())
}
