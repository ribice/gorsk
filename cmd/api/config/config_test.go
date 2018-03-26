package config_test

import (
	"reflect"
	"testing"

	"github.com/ribice/gorsk/cmd/api/config"
)

func TestLoad(t *testing.T) {
	type args struct {
		configName string
	}
	cases := []struct {
		name     string
		args     args
		wantData *config.Configuration
		wantErr  bool
	}{
		{
			name:    "Fail on non-existing file",
			args:    args{configName: "notExists"},
			wantErr: true,
		},
		{
			name:    "Fail on wrong file format",
			args:    args{configName: "invalid"},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{configName: "testdata"},
			wantData: &config.Configuration{
				DB: &config.DBConfig{
					Log:          true,
					CreateSchema: false,
				},
				Server: &config.ServerConfig{
					Port: 8080,
				},
				JWT: &config.JWTConfig{
					Timeout: 10800,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load(tt.args.configName)
			if !reflect.DeepEqual(tt.wantData, cfg) {
				t.Errorf("Expected and returned data does not match")
			}
			if tt.wantErr != (err != nil) {
				t.Errorf("Want err differs from err!=nil: %s", err)
			}
		})
	}
}
