package main

import (
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	//"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	gsessions "github.com/gorilla/sessions"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"

	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"go.uber.org/zap"

	bikeauth "gin-oauth/controllers"
	logger "gin-oauth/internal/logging"
)

var googleScopes = []string{
    "openid",
    "profile",
    "email",
    "https://www.googleapis.com/auth/drive.file",
    "https://www.googleapis.com/auth/drive.readonly",
    "https://www.googleapis.com/auth/drive",
}

var ilog logger.Logger
var loggerMgr *zap.Logger
var log *zap.SugaredLogger

const SessionName = "_gothic_session"
const SessionKey = "mysession"
var userStore gsessions.Store
var mydir string


func init() {
    workdir, err := os.Getwd() 
    mydir = workdir
    err = godotenv.Load(".env")
    if err != nil {
        log.Error(".env file failed to load!", err)
    }

    loggerMgr = logger.GetZapLogger()
    log = loggerMgr.Sugar()

    key := []byte(os.Getenv("SESSION_SECRET"))
    
	cookieStore := gsessions.NewCookieStore(key)
    cookieStore.Options = &gsessions.Options{
        Path:     "/",
        MaxAge:   3601,
        HttpOnly: true,
        Secure:   true,
    }
	cookieStore.Options.HttpOnly = true
	userStore = cookieStore

    gothic.Store = userStore

    clientID := os.Getenv("GOOGLE_CLIENT_ID")
    clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
    clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")
    if clientID =="" || clientSecret == "" || clientCallbackURL == "" {
        log.Error("Env vars GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET. CLIENT_CALLBACK_URL are required")
    }
    
    // Gin needs the information of the provider
    goth.UseProviders(
        google.New(clientID, clientSecret, clientCallbackURL, "email"),
    )

    log.Infof("Initialization complete")   
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

func main() {
    gin.SetMode(gin.ReleaseMode)
    r := gin.Default()

    bikeauth.InitRestHandler(log)
    
    r.Use(sessions.Sessions(SessionKey, NewCookieStore()))

    log.Infof("Setup session complete")   
    r.LoadHTMLGlob(mydir +"./templates/*")
    r.Static("/css", "./static/css")
    r.Static("/js", "./static/js")

    r.GET("/", bikeauth.Home)
    r.GET("/signup", bikeauth.Provider, bikeauth.SignUp)
    r.GET("/auth/:provider", bikeauth.Provider, bikeauth.SignInWithProvider)
    r.GET("/auth/:provider/callback", bikeauth.Provider, bikeauth.CallbackHandler)
    r.GET("/success", bikeauth.Success)
    r.GET("/auth/:provider/logout", bikeauth.Provider, bikeauth.Logout)
    log.Infof("+===================================================+")
    log.Infof("| Starting up the server on port 5000               |")   
    log.Infof("+===================================================+")

    r.Run(":5000")
}
