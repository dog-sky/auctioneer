package logging

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	type args struct {
		logLvl string
	}

	tests := []struct {
		name    string
		args    args
		want    *logrus.Logger
		wantErr bool
	}{
		{
			name: "OK INFO",
			args: args{
				"INFO",
			},
			want:    &logrus.Logger{},
			wantErr: false,
		},
		{
			name: "OK DEBUG",
			args: args{
				"DEBUG",
			},
			want:    &logrus.Logger{},
			wantErr: false,
		},
		{
			name: "OK ERROR",
			args: args{
				"ERROR",
			},
			want:    &logrus.Logger{},
			wantErr: false,
		},
		{
			name: "NOT OK",
			args: args{
				"NOT VALUE",
			},
			want:    &logrus.Logger{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewLogger(tt.args.logLvl)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)

				return
			}
			assert.IsType(t, tt.want, got)
		})
	}
}
