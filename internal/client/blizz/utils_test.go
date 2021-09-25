package blizz

import "testing"

func Test_isRussian(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "OK Russian",
			text: "Русский",
			want: true,
		}, {
			name: "OK English",
			text: "English",
			want: false,
		}, {
			name: "OK Russian with whitespace",
			text: "Русский с пробелами",
			want: true,
		}, {
			name: "OK English with whitespaces",
			text: "English whith whitespaces",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRussian(tt.text); got != tt.want {
				t.Errorf("isRussian() = %v, want %v", got, tt.want)
			}
		})
	}
}
