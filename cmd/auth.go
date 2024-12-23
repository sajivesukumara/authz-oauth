package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	"github.com/joho/godotenv"

	"go.uber.org/zap"

	bikeauth "gin-oauth/controllers"
	logger "gin-oauth/internal/logging"
	"gin-oauth/internal/utils"
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


func init() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Error(".env file failed to load!", err)
    }

    loggerMgr = logger.GetZapLogger()
    log = loggerMgr.Sugar()

    utils.SetupCookie()
    utils.SetupProviders(log)

    log.Infof("Initialization complete")   
}

func main() {
    gin.SetMode(gin.ReleaseMode)
    r := gin.Default()

    bikeauth.InitRestHandler(log)
    
    r.Use(sessions.Sessions(SessionKey, utils.NewCookieStore()))

    log.Infof("Setup session complete")   
    r.LoadHTMLGlob("./templates/*")
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
