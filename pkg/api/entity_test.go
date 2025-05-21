package api

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntityValidation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     Entity
		expectErr bool
		error     error
	}{
		{
			name:      "empty entity",
			input:     Entity{},
			expectErr: true,
		},
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validate.Struct(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.error, err)

				var validateErrs validator.ValidationErrors
				if errors.As(err, &validateErrs) {
					for _, e := range validateErrs {
						fmt.Println(e.Namespace())
						fmt.Println(e.Field())
						fmt.Println(e.StructNamespace())
						fmt.Println(e.StructField())
						fmt.Println(e.Tag())
						fmt.Println(e.ActualTag())
						fmt.Println(e.Kind())
						fmt.Println(e.Type())
						fmt.Println(e.Value())
						fmt.Println(e.Param())
						fmt.Println()
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
