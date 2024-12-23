package bikeauth

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
    "github.com/gin-contrib/sessions"
)


var log *zap.SugaredLogger

func InitRestHandler(logger *zap.SugaredLogger) {
    log = logger
}

// Load the home page using the template index.html file
// url : /
func Home(c *gin.Context) {
    log.Info("Loading home screen")

    tmpl, err := template.ParseFiles("./templates/index.html")
    if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
        log.Error(err)
        return
    }
    err = tmpl.Execute(c.Writer, gin.H{})
    if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    log.Info("Loaded home page")
}

// Load the home page using the template index.html file
// url : /
func SignUp(c *gin.Context) {
    log.Info("Loading Signup Page")
    tmpl, err := template.ParseFiles("templates/signup.html")
    if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
        log.Error(err)
        return
    }
    err = tmpl.Execute(c.Writer, gin.H{})
    if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    log.Info("Loaded home page")
}


// url : auth/:provider
func SignInWithProvider(c *gin.Context) {
    log.Infof("Redirect URL : %s", DefineUrl(c))
    user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
    if err != nil {
        gothic.BeginAuthHandler(c.Writer, c.Request)
        log.Debugf("[%s] : Sign in as new user [%v]", c.Param("provider"), err)
        
        return
    }
	
    profile := NewProfile(user)
	session := sessions.Default(c)
	session.Set("profile", profile)

    session1 := sessions.Default(c)
    session1.Set("FirstName", "sajivekumar")

    log.Debugf("[%s] : Sign in existing user : %s", c.Param("provider"), user)
}

// url : auth/:provider/callback
func CallbackHandler(c *gin.Context){
    userInfo ,err := gothic.CompleteUserAuth(c.Writer, c.Request)
    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }
    profile := NewProfile(userInfo)
	session := sessions.Default(c)
	session.Set("profile", profile)
    session.Save()

    log.Infof("Cookies Access : %s", userInfo.AccessToken)
    log.Infof("Cookies Refresh: %s", userInfo.RefreshToken)
    userResponse := LoginResponse { userInfo.AccessToken, userInfo.RefreshToken}
    SetCookie(c, userResponse)
    
    log.Infof("Response User info : %+v", userInfo)
    c.Redirect(http.StatusTemporaryRedirect, "/success")
}

// url : /success
func Success(c *gin.Context) {
   
    cookie := GetCookie(c)
    if cookie == "" {
        c.Redirect(http.StatusTemporaryRedirect, "/")
    } else {
        c.JSON(200, cookie)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf(`
      <div style="
          background-color: #fff;
          padding: 40px;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          text-align: center;
      ">
          <h1 style="
              color: #333;
              margin-bottom: 20px;
          ">You have Successfull signed in!</h1>
          
          </div>
      </div>
  `)))
    }
}

// url : /auth/logout/:provider
func Logout(c *gin.Context) {
    gothic.Logout(c.Writer, c.Request)
    SetCookie(c, LoginResponse{})
    log.Info("Cleaned up cookies")
    c.Writer.Header().Set("Location", "/")
    c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}
