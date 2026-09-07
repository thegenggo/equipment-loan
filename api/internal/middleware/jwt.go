package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thegenggo/equipment-loan/api/internal/httperr"
	"github.com/thegenggo/equipment-loan/api/pkg/token"
)

const (
	ContextUserID = "user_id"
	ContextRole   = "role"
)

func Auth(tokens *token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, found := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !found || raw == "" {
			httperr.Write(c, http.StatusUnauthorized, "unauthorized", "missing bearer token")
			return
		}

		claims, err := tokens.Verify(raw)
		if err != nil {
			httperr.Write(c, http.StatusUnauthorized, "unauthorized", "invalid or expired token")
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)

		c.Next()
	}
}
