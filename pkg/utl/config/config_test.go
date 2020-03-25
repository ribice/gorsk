package config_test

import (
	"testing"

	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		wantData *config.Configuration
		wantErr  bool
	}{
		{
			name:    "Fail on non-existing file",
			path:    "notExists",
			wantErr: true,
		},
		{
			name:    "Fail on wrong file format",
			path:    "testdata/config.invalid.yaml",
			wantErr: true,
		},
		{
			name: "Success",
			path: "testdata/config.testdata.yaml",
			wantData: &config.Configuration{
				DB: &config.Database{
					LogQueries: true,
					Timeout:    20,
				},
				Server: &config.Server{
					Port:         ":8080",
					Debug:        true,
					ReadTimeout:  15,
					WriteTimeout: 20,
				},
				JWT: &config.JWT{
					MinSecretLength:  128,
					DurationMinutes:  10,
					RefreshDuration:  10,
					MaxRefresh:       144,
					SigningAlgorithm: "HS384",
				},
				App: &config.Application{
					MinPasswordStr: 3,
					SwaggerUIPath:  "assets/swagger",
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load(tt.path)
			assert.Equal(t, tt.wantData, cfg)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
