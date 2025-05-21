package api

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func testEntityValidation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     *Entity
		expectErr bool
		errors    []struct {
			Namespace string
			Tag       string
		}
	}{
		{
			name: "empty entity",
			input: &Entity{
				Name:         "",
				Type:         "",
				Observations: nil,
			},
			expectErr: true,
			errors: []struct {
				Namespace string
				Tag       string
			}{
				{
					Namespace: "Entity.Type",
					Tag:       "required",
				},
				{
					Namespace: "Entity.Name",
					Tag:       "required",
				},
			},
		},
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validate.Struct(tt.input)
			if tt.expectErr {
				assert.Error(t, err)

				var validateErrs validator.ValidationErrors
				if errors.As(err, &validateErrs) {
					assert.Equal(t, 0, len(validateErrs))
					for _, e := range validateErrs {
						assert.Equal(t, tt.errors[i].Namespace, e.Namespace())
						assert.Equal(t, tt.errors[i].Tag, e.Tag())
						t.Log(e.Namespace())
						t.Log(e.Field())
						t.Log(e.StructNamespace())
						t.Log(e.StructField())
						t.Log(e.Tag())
						t.Log(e.ActualTag())
						t.Log(e.Kind())
						t.Log(e.Type())
						t.Log(e.Value())
						t.Log(e.Param())
						t.Log()
					}
				} else {
					assert.Fail(t, "Expected validation errors, but got: %T", err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
