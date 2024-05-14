package usersession

type UserSession struct {
	Username string

	Passwords []struct {
		platform string
		password string
	}
}

func NewUserSession(username string, passwords []struct{ platform, password string }) *UserSession {
	return &UserSession{
		Username:  username,
		Passwords: passwords,
	}
}
