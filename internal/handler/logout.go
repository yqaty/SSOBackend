package handler

import (
	"github.com/UniqueStudio/UniqueSSOBackend/internal/constants"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Set(constants.SessionNameUID, "")
	sess.Options(sessions.Options{MaxAge: -1})
	sess.Save()
}
