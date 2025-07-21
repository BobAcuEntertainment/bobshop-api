package security

type Tokenizer interface {
	GenerateToken(userID string, role string) (string, error)
	ParseToken(token string) (map[string]any, error)
}
