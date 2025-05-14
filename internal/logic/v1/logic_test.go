package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type unmarshalable struct{}

func (u unmarshalable) MarshalJSON() ([]byte, error) {
	return nil, errors.New("cannot marshal")
}

func Test_toolJSONResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		wantErr  bool
		wantJSON string
	}{
		{
			name:    "simple map",
			input:   map[string]int{"foo": 1},
			wantErr: false,
			wantJSON: `{
  "foo": 1
}`,
		},
		{
			name:    "unmarshalable type",
			input:   unmarshalable{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := toolJSONResponse(context.Background(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				if assert.NotNil(t, resp) {
					if assert.Len(t, resp.Content, 1) {
						assert.Equal(t, resp.Content[0].TextContent.Text, tt.wantJSON)
					}
				}
			}
		})
	}
}
