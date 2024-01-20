package main

import (
	"fmt"
	"os"
	"github.com/goplus/yap"
	"net/http"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"go.uber.org/zap"
	_ "github.com/joho/godotenv/autoload"
	"github.com/goplus/account/internal/core"
)

type account struct {
	yap.App
}
//line cmd/account_yap.gox:12
func (this *account) MainEntry() {
//line cmd/account_yap.gox:12:1
	endpoint := os.Getenv("GOP_ACCOUNT_ENDPOINT")
//line cmd/account_yap.gox:13:1
	logger, _ := zap.NewProduction()
//line cmd/account_yap.gox:14:1
	defer logger.Sync()
//line cmd/account_yap.gox:15:1
	zlog := logger.Sugar()
//line cmd/account_yap.gox:17:1
	this.Get("/p/:id", func(ctx *yap.Context) {
//line cmd/account_yap.gox:18:1
		ctx.Yap__1("article", map[string]string{"id": ctx.Param("id")})
	})
//line cmd/account_yap.gox:22:1
	this.Get("/", func(ctx *yap.Context) {
//line cmd/account_yap.gox:23:1
		ctx.Yap__1("home", map[string]interface {
		}{})
	})
//line cmd/account_yap.gox:26:1
	this.Get("/login", func(ctx *yap.Context) {
//line cmd/account_yap.gox:27:1
		casdoorURL := "https://casdoor-community.qiniu.io"
//line cmd/account_yap.gox:28:1
		clientID := "49a8ac9729e314a05bf0"
//line cmd/account_yap.gox:29:1
		redirectURI := "http://localhost:8081/callback"
//line cmd/account_yap.gox:30:1
		ResponseType := "code"
//line cmd/account_yap.gox:31:1
		Scope := "read"
//line cmd/account_yap.gox:32:1
		State := "casdoor"
//line cmd/account_yap.gox:33:1
		loginURL := fmt.Sprintf("%s/login/oauth/authorize?client_id=%s&response_type=%s&redirect_uri=%s&scope=%s&state=%s", casdoorURL, clientID, ResponseType, redirectURI, Scope, State)
//line cmd/account_yap.gox:34:1
		http.Redirect(ctx.ResponseWriter, ctx.Request, loginURL, http.StatusFound)
	})
//line cmd/account_yap.gox:37:1
	this.Get("/callback", func(ctx *yap.Context) {
//line cmd/account_yap.gox:38:1
		code := ctx.URL.Query().Get("code")
//line cmd/account_yap.gox:39:1
		state := ctx.URL.Query().Get("state")
//line cmd/account_yap.gox:40:1
		fmt.Println("--code and state--")
//line cmd/account_yap.gox:41:1
		fmt.Println("code:", code)
//line cmd/account_yap.gox:42:1
		fmt.Println("state:", state)
//line cmd/account_yap.gox:43:1
		token, err := casdoorsdk.GetOAuthToken(code, state)
//line cmd/account_yap.gox:44:1
		if err != nil {
//line cmd/account_yap.gox:45:1
			fmt.Println("err", err)
		}
//line cmd/account_yap.gox:47:1
		claim, err := casdoorsdk.ParseJwtToken(token.AccessToken)
//line cmd/account_yap.gox:48:1
		if err != nil {
//line cmd/account_yap.gox:49:1
			zlog.Error("err", err)
		}
//line cmd/account_yap.gox:51:1
		username := claim.User.Name
//line cmd/account_yap.gox:53:1
		cookie := http.Cookie{Name: "token", Value: token.AccessToken, Path: "/", MaxAge: 3600}
//line cmd/account_yap.gox:59:1
		http.SetCookie(ctx.ResponseWriter, &cookie)
//line cmd/account_yap.gox:61:1
		ctx.Yap__1("callback", map[string]string{"username": username, "usertoken": token.AccessToken})
	})
//line cmd/account_yap.gox:67:1
	this.Get("/test", func(ctx *yap.Context) {
//line cmd/account_yap.gox:68:1
		cookie, err := ctx.Request.Cookie("token")
//line cmd/account_yap.gox:69:1
		if err != nil {
//line cmd/account_yap.gox:70:1
			zlog.Error("err", err)
		}
//line cmd/account_yap.gox:72:1
		claim, err := casdoorsdk.ParseJwtToken(cookie.Value)
//line cmd/account_yap.gox:73:1
		if err != nil {
//line cmd/account_yap.gox:74:1
			zlog.Error("err", err)
		} else {
//line cmd/account_yap.gox:76:1
			fmt.Println("username:", claim.User.Name)
//line cmd/account_yap.gox:77:1
			fmt.Println("userId:", claim.User.Id)
//line cmd/account_yap.gox:78:1
			fmt.Println("password:", claim.User.Password)
		}
//line cmd/account_yap.gox:80:1
		ctx.Yap__1("test", map[string]interface {
		}{})
	})
//line cmd/account_yap.gox:83:1
	core.Init()
//line cmd/account_yap.gox:84:1
	zlog.Info("Started in endpoint: ", endpoint)
//line cmd/account_yap.gox:85:1
	this.Run(endpoint)
}
func main() {
	yap.Gopt_App_Main(new(account))
}
