package handler

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/UniqueStudio/UniqueSSOBackend/internal/constants"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/tracer"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

type sessionUserUID struct{}
type sessionUserPermission struct{}

func mustGetUIDFromCtx(ctx context.Context) string {
	return ctx.Value(sessionUserUID{}).(string)
}

func ctxWithUID(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, sessionUserUID{}, uid)
}

// AuthenticationMiddleware - for authentication. check by sessionID
// if not auth, it will redirect to ${ssoLoginURL}?${constants.SSORedirectBackQueryKey}=${currentURL}
func AuthenticationMiddleware(c *gin.Context) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "Authentication")
	defer span.End()

	// Add byPass cookie
	cookie, _ := c.Cookie("SSO_SESSION")
	if cookie == "unique_web_admin" {
		c.Request = c.Request.WithContext(ctxWithUID(apmCtx, "ffb6e834-3615-4ebb-9d9d-825af333a3ca"))
		span.SetAttributes(attribute.String("UID", "ffb6e834-3615-4ebb-9d9d-825af333a3ca"))
		c.Next()
		return
	}

	s := sessions.Default(c)
	u := s.Get(constants.SessionNameUID)

	if u == nil {
		c.Abort()
		respForbiddenError(c, errors.New("unauthenticated"))
		return
	}
	uid, ok := u.(string)
	if !ok {
		c.Abort()
		respForbiddenError(c, errors.New("unauthenticated"))
		return
	}
	c.Request = c.Request.WithContext(ctxWithUID(apmCtx, uid))

	// TODO(xylonx): auto renew session

	span.SetAttributes(attribute.String("UID", uid))
}

func RedirectMiddleware(c *gin.Context) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "Authentication")
	defer span.End()

	// Add byPass cookie
	cookie, _ := c.Cookie("SSO_SESSION")
	if cookie == "unique_web_admin" {
		c.Request = c.Request.WithContext(ctxWithUID(apmCtx, "ffb6e834-3615-4ebb-9d9d-825af333a3ca"))
		span.SetAttributes(attribute.String("UID", "ffb6e834-3615-4ebb-9d9d-825af333a3ca"))
		c.Next()
		return
	}

	s := sessions.Default(c)
	u := s.Get(constants.SessionNameUID)

	if u == nil {
		c.Abort()
		c.Redirect(http.StatusFound, "http://81.70.253.156:54250")
		return
	}
	uid, ok := u.(string)
	if !ok {
		c.Abort()
		c.Redirect(http.StatusFound, "http://81.70.253.156:54250")
		return
	}
	c.Request = c.Request.WithContext(ctxWithUID(apmCtx, uid))

	// TODO(xylonx): auto renew session

	span.SetAttributes(attribute.String("UID", uid))
}

func TraefikAuth(forwardAuthURI *url.URL) gin.HandlerFunc {
	return func(c *gin.Context) {
		apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "TraefikForwardAuth")
		defer span.End()

		// method := c.GetHeader("X-Forwarded-Method")
		scheme := c.GetHeader("X-Forwarded-Proto")
		host := c.GetHeader("X-Forwarded-Host")
		uri := c.GetHeader("X-Forwarded-Uri")
		originURI := scheme + ":" + host + uri

		s := sessions.Default(c)
		uid := s.Get(constants.SessionNameUID)
		// Add byPass code
		// due to recruitment and hackday system are under development
		// so pass the request to https://hr.hustunique.com/ and https://join.hustunique.com/
		// and https://hackday2023.hustunique.com/
		// the forwardAuth api is https://SSOURL/gateway/validate/traefik

		if host == "hr.hustunique.com" || host == "join.hustunique.com" || host == "hackday2023.hustunique.com" || host == "back.recruitment.hustunique.com" {
			c.Request = c.Request.WithContext(ctxWithUID(apmCtx, uid.(string)))
			respOK(c, "success")
			return
		}
		uidString, ok := uid.(string)
		if uid == nil || !ok || uidString != "" {
			forwardURI := *forwardAuthURI
			q := forwardURI.Query()
			q.Set(constants.SSORedirectBackQueryKey, originURI)
			forwardURI.RawQuery = q.Encode()

			c.Abort()
			c.Redirect(http.StatusTemporaryRedirect, forwardURI.String())
			return
		}

		// add X-UID for success http request
		c.Header("X-UID", uid.(string))

		c.Request = c.Request.WithContext(ctxWithUID(apmCtx, uid.(string)))
		respOK(c, "success")
	}
}
