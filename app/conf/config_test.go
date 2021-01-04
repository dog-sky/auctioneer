package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{
			name:   "ok",
			want:   &Config{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig()
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.IsType(t, tt.want, got)
		})
	}
}
