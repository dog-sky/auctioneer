package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name       string
		want       *Config
		setGoodEnv bool
		wantErr    bool
	}{
		{
			name:       "not ok",
			want:       &Config{},
			setGoodEnv: false,
			wantErr:    true,
		},
		{
			name:       "ok",
			want:       &Config{},
			setGoodEnv: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setGoodEnv {
				setGoodEnv()
			}
			got, err := NewConfig()
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.IsType(t, tt.want, got)
		})
	}
}

func setGoodEnv() {
	os.Setenv("AUCTIONEER_BLIZZARD_CLIENT_SECRET", "MattDaemon")
	os.Setenv("AUCTIONEER_BLIZZARD_CLIENT_ID", "MattDaemon")
	os.Setenv("AUCTIONEER_APP_PORT", ":8000")
}
