package model

type CoreUser interface {
	GetID() string
	GetEmail() string
	GetPasswordHash() string
}

type User struct {
	ID       string
	Email    string
	Password string
}

func (u User) GetID() string           { return u.ID }
func (u User) GetEmail() string        { return u.Email }
func (u User) GetPasswordHash() string { return u.Password }
