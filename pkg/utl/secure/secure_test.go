package secure_test

import (
	"crypto/sha1"
	"testing"

	"github.com/ribice/gorsk/pkg/utl/secure"
	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	cases := []struct {
		name   string
		pass   string
		inputs []string
		want   bool
	}{
		{
			name: "Insecure password",
			pass: "notSec",
			want: false,
		},
		{
			name:   "Password matches input fields",
			pass:   "johndoe92",
			inputs: []string{"John", "Doe"},
			want:   false,
		},
		{
			name:   "Secure password",
			pass:   "callgophers",
			inputs: []string{"John", "Doe"},
			want:   true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := secure.New(1, nil)
			got := s.Password(tt.pass, tt.inputs...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHashAndMatch(t *testing.T) {
	cases := []struct {
		name string
		pass string
		want bool
	}{
		{
			name: "Success",
			pass: "gamepad",
			want: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := secure.New(1, nil)
			hash := s.Hash(tt.pass)
			assert.Equal(t, tt.want, s.HashMatchesPassword(hash, tt.pass))
		})
	}
}

func TestToken(t *testing.T) {
	s := secure.New(1, sha1.New())
	token := "token"
	tokenized := s.Token(token)
	assert.NotEqual(t, tokenized, token)
}
