package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thegenggo/equipment-loan/api/internal/httperr"
)

func pathID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		httperr.Write(c, http.StatusBadRequest, "invalid_id", "id must be a positive integer")
		return 0, false
	}

	return id, true
}
