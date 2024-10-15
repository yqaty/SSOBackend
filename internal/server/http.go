package server

import (
	"net/http"
	"time"

	"github.com/UniqueStudio/UniqueSSOBackend/internal/core"
	"github.com/gin-contrib/sessions"
	"github.com/xylonx/zapx"

	"github.com/UniqueStudio/UniqueSSOBackend/config"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/handler"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/tracer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HttpOption struct {
	Addr         string
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	AllowOrigins []string
}

func InitHttpServer(o *HttpOption) *http.Server {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(tracer.TracingMiddleware)
	gin.SetMode(o.Mode)

	// config cors
	if o.Mode == "debug" {
		zapx.Infof("service is in debug mode")
		corsConfig := cors.DefaultConfig()
		//corsConfig.AllowAllOrigins = true
		corsConfig.AllowOrigins = o.AllowOrigins
		corsConfig.AllowCredentials = true
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
		r.Use(cors.New(corsConfig))
	} else {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = o.AllowOrigins
		corsConfig.AllowCredentials = true
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
		r.Use(cors.New(corsConfig))
	}

	r.Use(sessions.Sessions("SSO_SESSION", core.SessStore))

	// ping router
	r.GET("/ping", handler.RedirectMiddleware, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "this is uniquestudio sso system",
		})
	})

	r.GET("/api/v1/ping", handler.RedirectMiddleware, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "this is uniquestudio sso system",
		})
	})

	router := r.Group("/api/v1")

	// register router
	registerSub := router.Group("/register")
	registerSub.POST("", handler.RegisterHandler)

	resetSub := router.Group("/reset")
	resetSub.POST("", handler.ResetPassword)

	// login router
	loginSub := router.Group("/login")
	loginSub.POST("", handler.LoginHandler)

	// logout router
	logoutSub := router.Group("/logout")
	logoutSub.POST("", handler.LogoutHandler)

	// 2fa code router
	codeSub := router.Group("/code")
	codeSub.POST("/sms", handler.SendSMSCodeHandler)
	//codeSub.POST("/email", handler.SendEmailCodeHandler)

	// user router
	userSub := router.Group("/user")
	userSub.Use(handler.AuthenticationMiddleware)
	userSub.GET("/my", handler.GetUserInfoHandler)
	userSub.PUT("/my", handler.UpdateUserInfoHandler)

	// traefik redirect router
	gatewaySub := router.Group("/gateway")
	gatewaySub.GET("/validate/traefik", handler.TraefikAuth(config.TraefikRedirectURI))

	return &http.Server{
		Addr:         o.Addr,
		Handler:      r,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,
	}
}
