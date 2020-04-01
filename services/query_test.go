package services

import (
	"reflect"
	"testing"

	"github.com/logiqai/logiq-box/api/v1/query"
)

func Test_getFilters(t *testing.T) {
	type args struct {
		filters string
	}
	tests := []struct {
		name string
		args args
		want map[string]*query.FilterValues
	}{
		{
			name: "single",
			args: args{"HostName=127.0.0.1"},
			want: map[string]*query.FilterValues{
				"HostName": {
					Values: []string{"127.0.0.1"},
				},
			},
		},
		{
			name: "single-with-multi",
			args: args{"HostName=127.0.0.1,logiq.ai"},
			want: map[string]*query.FilterValues{
				"HostName": {
					Values: []string{"127.0.0.1", "logiq.ai"},
				},
			},
		},
		{
			name: "multi-multi",
			args: args{"HostName=127.0.0.1,logiq.ai;Message=ghost,tar,jar"},
			want: map[string]*query.FilterValues{
				"HostName": {
					Values: []string{"127.0.0.1", "logiq.ai"},
				},
				"Message": {
					Values: []string{"ghost", "tar", "jar"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFilters(tt.args.filters); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}
