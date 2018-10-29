package mock

// Secure mock
type Secure struct {
	PasswordFn            func(string, ...string) bool
	HashFn                func(string) string
	HashMatchesPasswordFn func(string, string) bool
	TokenFn               func(string) string
}

// Password mock
func (s *Secure) Password(pw string, inputs ...string) bool {
	return s.PasswordFn(pw, inputs...)
}

// Hash mock
func (s *Secure) Hash(pw string) string {
	return s.HashFn(pw)
}

// HashMatchesPassword mock
func (s *Secure) HashMatchesPassword(hash, pw string) bool {
	return s.HashMatchesPasswordFn(hash, pw)
}

// Token mock
func (s *Secure) Token(token string) string {
	return s.TokenFn(token)
}
