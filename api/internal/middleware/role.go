package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thegenggo/equipment-loan/api/internal/httperr"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		value, exists := c.Get(ContextRole)
		role, ok := value.(string)
		if !exists || !ok {
			httperr.Write(c, http.StatusUnauthorized, "unauthorized", "authentication required")
			return
		}

		if _, ok := allowed[role]; !ok {
			httperr.Write(c, http.StatusForbidden, "forbidden", "insufficient permission")
			return
		}

		c.Next()
	}
}
