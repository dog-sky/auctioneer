package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		logLvl string
	}
	tests := []struct {
		name    string
		args    args
		want    *Logger
		wantErr bool
	}{
		{
			name: "OK INFO",
			args: args{
				"INFO",
			},
			want:    &Logger{},
			wantErr: false,
		},
		{
			name: "OK DEBUG",
			args: args{
				"DEBUG",
			},
			want:    &Logger{},
			wantErr: false,
		},
		{
			name: "OK ERROR",
			args: args{
				"ERROR",
			},
			want:    &Logger{},
			wantErr: false,
		},
		{
			name: "NOT OK",
			args: args{
				"NOT VALUE",
			},
			want:    &Logger{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLogger(tt.args.logLvl)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.IsType(t, tt.want, got)
		})
	}
}
