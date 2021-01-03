package v1

import "testing"

func Test_checkQueryParams(t *testing.T) {
	type args struct {
		q *searchQueryParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK all params passed",
			args: args{
				q: &searchQueryParams{
					RealmName: "hello",
					ItemName:  "hello",
					Region:    "hello",
				},
			},
			wantErr: false,
		},
		{
			name: "ERR No RealmName",
			args: args{
				q: &searchQueryParams{
					ItemName: "hello",
					Region:   "hello",
				},
			},
			wantErr: true,
		},
		{
			name: "ERR No Region",
			args: args{
				q: &searchQueryParams{
					RealmName: "hello",
					ItemName:  "hello",
				},
			},
			wantErr: true,
		},
		{
			name: "ERR No ItemName",
			args: args{
				q: &searchQueryParams{
					RealmName: "hello",
					Region:    "hello",
				},
			},
			wantErr: true,
		},
		{
			name: "ERR No params",
			args: args{
				q: &searchQueryParams{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkQueryParams(tt.args.q); (err != nil) != tt.wantErr {
				t.Errorf("checkQueryParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
