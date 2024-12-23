package bikeauth

import (
	"time"

	"github.com/markbates/goth"
)

type Profile struct {
	ID           string    `json:"id"`
	AvatarURL    string    `json:"avatar"`
	NickName     string    `json:"nickname"`
	Lastname     string    `json:"lastname"`
	Firstname    string    `json:"firstname"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Groups       []string  `json:"groups"`
	RefreshToken string    `json:"-"`
	AccessToken  string    `json:"-"`
	IDToken      string    `json:"-"`
	ExpiresAt    time.Time `json:"-"`
}

func NewProfile(user goth.User) Profile {
	return Profile{
		ID:           user.UserID,
		AvatarURL:    user.AvatarURL,
		NickName:     user.NickName,
		Firstname:    user.FirstName,
		Lastname:     user.LastName,
		Name:         user.Name,
		Email:        user.Email,
		ExpiresAt:    user.ExpiresAt,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		IDToken:      user.IDToken,
	}
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}