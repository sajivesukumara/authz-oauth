package bikeauth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

const (
    ProviderKey = "provider"
)

func DefineUrl(c *gin.Context) string {
    provider := c.Param("provider")
    q := c.Request.URL.Query()
    q.Add("provider", provider)
    q.Add("state", gothic.SetState(c.Request))
    c.Request.URL.RawQuery = q.Encode()
    url, err := gothic.GetAuthURL(c.Writer, c.Request)
    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
    }
    return url
}

func Provider(ctx *gin.Context) {
    provider := ctx.Param("provider")
    log.Infof("Provider middleware triggered %s", provider)
    ctx.Request = ctx.Request.WithContext(context.WithValue(ctx, ProviderKey, provider))
    ctx.Set(ProviderKey, provider)
    ctx.Next()
}

func SetCookie(c *gin.Context, tokens LoginResponse) error {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	/* session, err := gothic.Store.Get(c.Request(), s.cfg.SessionName()) */
	/* if err != nil { */
	/* 	logrus.Errorf("StoreUserSession.Get(): %v\n", err) */
	/* 	return err */
	/* } */

	/* session.Values["authentication"] = string(tokensJSON) */
    if tokens == (LoginResponse{}) {
        c.SetCookie("bike-access-token",  "", 0, "", "", true, true)
        c.SetCookie("bike-refresh-token",  "", 0, "", "", true, true)
    }

	c.SetCookie("bike-access-token",  tokens.AccessToken, 3600*24, "", "", true, true)
    c.SetCookie("bike-refresh-token",  tokens.RefreshToken, 3600*24, "", "", true, true)

	return nil
}

func GetCookie(c *gin.Context) string {
    cookie, err := c.Cookie("bike-access-token")
    if err != nil {
        log.Error("failed to read cookie", err)
        return ""
    }
    log.Infof("Found cookie : %s", cookie)
    return cookie
}