package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

type Config struct {
    GoogleLoginConfig oauth2.Config
}

var AppConfig Config

func GoogleConfig() oauth2.Config {
    err := godotenv.Load("gin-oauth/.env")
    if err != nil {
        log.Fatal(".env file failed to load!")
    }

    log.Printf("Read env file %s", os.Getenv("GOOGLE_CLIENT_ID"))

    AppConfig.GoogleLoginConfig = oauth2.Config{
        RedirectURL:  os.Getenv("CLIENT_CALLBACK_URL"),
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile"},
        Endpoint: google.Endpoint,
    }

    return AppConfig.GoogleLoginConfig
}

