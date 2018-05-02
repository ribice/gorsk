package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
				DB: &config.Database{
					Log:          true,
					CreateSchema: false,
				},
				Server: &config.Server{
					Port:  ":8080",
					Debug: true,
				},
				JWT: &config.JWT{
					Duration: 10800,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load(tt.args.configName)
			assert.Equal(t, tt.wantData, cfg)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
