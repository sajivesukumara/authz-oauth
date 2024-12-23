package utils

import (
	"errors"
	"os"

	//logger "gin-oauth/internal/logging"
	gsessions "github.com/gorilla/sessions"
	"go.uber.org/zap"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var userStore gsessions.Store

func SetupProviders(log *zap.SugaredLogger){ 
    clientID := os.Getenv("GOOGLE_CLIENT_ID")
    clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
    clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")
    if clientID =="" || clientSecret == "" || clientCallbackURL == "" {
        log.Error("Env vars GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET. CLIENT_CALLBACK_URL are required",  errors.New("Env missing"))
    }
    
    // Gin needs the information of the provider
    goth.UseProviders(
        google.New(clientID, clientSecret, clientCallbackURL, "email"),
    )
}

func SetupCookie(){
    session_secret := os.Getenv("SESSION_SECRET")
    if session_secret == "" {
        session_secret = "my-default-secret"
    }
    
	cookieStore := gsessions.NewCookieStore([]byte(session_secret))
    cookieStore.Options = &gsessions.Options{
        Path:     "/",
        MaxAge:   3601,
        HttpOnly: true,
        Secure:   true,
    }
	cookieStore.Options.HttpOnly = true
	userStore = cookieStore

    gothic.Store = userStore
}

type store struct {
	gsessions.Store
}
func (c *store) Options(_ sessions.Options) {}

func NewStore(s gsessions.Store) sessions.Store {
	return &store{s}
}

func NewCookieStore() cookie.Store {
    localStore := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
    return localStore
}