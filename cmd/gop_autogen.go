package main

import (
	"os"
	"github.com/goplus/yap"
	"net/http"
	"go.uber.org/zap"
	"github.com/goplus/account/internal/core"
	_ "github.com/joho/godotenv/autoload"
)

type account struct {
	yap.App
	account *core.Account
}
//line cmd/account_yap.gox:16
func (this *account) MainEntry() {
//line cmd/account_yap.gox:16:1
	endpoint := os.Getenv("GOP_ACCOUNT_ENDPOINT")
//line cmd/account_yap.gox:17:1
	logger, _ := zap.NewProduction()
//line cmd/account_yap.gox:18:1
	defer logger.Sync()
//line cmd/account_yap.gox:19:1
	zlog := logger.Sugar()
//line cmd/account_yap.gox:21:1
	this.Get("/", func(ctx *yap.Context) {
//line cmd/account_yap.gox:22:1
		ctx.Yap__1("home", map[string]interface {
		}{})
	})
//line cmd/account_yap.gox:25:1
	this.Get("/login", func(ctx *yap.Context) {
//line cmd/account_yap.gox:27:1
		redirectURL := ctx.URL.Query().Get("redirect_url")
//line cmd/account_yap.gox:28:1
		loginURL := this.account.RedirectToCasdoor(redirectURL)
//line cmd/account_yap.gox:29:1
		ctx.Redirect(loginURL, http.StatusFound)
	})
//line cmd/account_yap.gox:32:1
	this.Get("/callback", func(ctx *yap.Context) {
//line cmd/account_yap.gox:33:1
		code := ctx.URL.Query().Get("code")
//line cmd/account_yap.gox:34:1
		state := ctx.URL.Query().Get("state")
//line cmd/account_yap.gox:36:1
		token, error := this.account.GetAccessToken(code, state)
//line cmd/account_yap.gox:37:1
		if error != nil {
//line cmd/account_yap.gox:38:1
			zlog.Error("err", error)
		}
//line cmd/account_yap.gox:41:1
		cookie := http.Cookie{Name: "token", Value: token.AccessToken, Path: "/", MaxAge: 3600}
//line cmd/account_yap.gox:47:1
		http.SetCookie(ctx.ResponseWriter, &cookie)
//line cmd/account_yap.gox:51:1
		http.Redirect(ctx.ResponseWriter, ctx.Request, "http://localhost:8080", http.StatusFound)
	})
//line cmd/account_yap.gox:54:1
	this.account = core.New()
//line cmd/account_yap.gox:55:1
	zlog.Info("Started in endpoint: ", endpoint)
//line cmd/account_yap.gox:56:1
	this.Run(endpoint)
}
func main() {
	yap.Gopt_App_Main(new(account))
}
