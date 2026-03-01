package ports

type Client struct {
	UserID  string
	ID      string
	Name    string
	Surname string
}

type TokenPair struct {
	AccessToken      string
	AccessExpiresAt  int64
	RefreshToken     string
	RefreshExpiresAt int64
	SessionId        string
}
