package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thegenggo/equipment-loan/api/internal/httperr"
	"github.com/thegenggo/equipment-loan/api/internal/middleware"
	"github.com/thegenggo/equipment-loan/api/internal/model"
	"github.com/thegenggo/equipment-loan/api/internal/service"
)

type createLoanRequest struct {
	EquipmentID int64  `json:"equipment_id" binding:"required,gt=0"`
	Purpose     string `json:"purpose" binding:"required,max=500"`
}

type LoanHandler struct {
	loans *service.LoanService
}

func NewLoanHandler(loans *service.LoanService) *LoanHandler {
	return &LoanHandler{loans: loans}
}

func (h *LoanHandler) Create(c *gin.Context) {
	var req createLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.Write(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	loan, err := h.loans.Create(c.Request.Context(), actorID(c), req.EquipmentID, req.Purpose)
	if err != nil {
		h.writeError(c, err, "create loan")
		return
	}

	c.JSON(http.StatusCreated, loan)
}

func (h *LoanHandler) List(c *gin.Context) {
	loans, err := h.loans.List(c.Request.Context(), actorID(c), actorRole(c), c.Query("status"))
	if err != nil {
		h.writeError(c, err, "list loans")
		return
	}

	c.JSON(http.StatusOK, loans)
}

func (h *LoanHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}

	loan, err := h.loans.Get(c.Request.Context(), id, actorID(c), actorRole(c))
	if err != nil {
		h.writeError(c, err, "get loan")
		return
	}

	c.JSON(http.StatusOK, loan)
}

func (h *LoanHandler) Approve(c *gin.Context) {
	h.decide(c, h.loans.Approve, "approve loan")
}

func (h *LoanHandler) Reject(c *gin.Context) {
	h.decide(c, h.loans.Reject, "reject loan")
}

func (h *LoanHandler) Return(c *gin.Context) {
	h.decide(c, h.loans.Return, "return loan")
}

func (h *LoanHandler) decide(
	c *gin.Context,
	action func(ctx context.Context, id, actorID int64) (*model.LoanRequestDetail, error),
	operation string,
) {
	id, ok := pathID(c)
	if !ok {
		return
	}

	loan, err := action(c.Request.Context(), id, actorID(c))
	if err != nil {
		h.writeError(c, err, operation)
		return
	}

	c.JSON(http.StatusOK, loan)
}

func (h *LoanHandler) writeError(c *gin.Context, err error, operation string) {
	switch {
	case errors.Is(err, service.ErrLoanNotFound):
		httperr.Write(c, http.StatusNotFound, "not_found", "loan request not found")
	case errors.Is(err, service.ErrEquipmentNotFound):
		httperr.Write(c, http.StatusNotFound, "not_found", "equipment not found")
	case errors.Is(err, service.ErrNotOwner):
		httperr.Write(c, http.StatusForbidden, "forbidden", "this request belongs to someone else")
	case errors.Is(err, service.ErrEquipmentNotAvailable):
		httperr.Write(c, http.StatusConflict, "equipment_unavailable", "this equipment is not available to borrow")
	case errors.Is(err, service.ErrEquipmentRequested):
		httperr.Write(c, http.StatusConflict, "equipment_requested", "this equipment already has an open request")
	case errors.Is(err, service.ErrLoanNotPending):
		httperr.Write(c, http.StatusConflict, "not_pending", "this request has already been decideed")
	case errors.Is(err, service.ErrLoanNotApproved):
		httperr.Write(c, http.StatusConflict, "not_on_loan", "only an approved request can be returned")
	case errors.Is(err, service.ErrInvalidStatus):
		httperr.Write(c, http.StatusBadRequest, "invalid_status", "unknown loan status")
	default:
		log.Printf("%s: %v", operation, err)
		httperr.Write(c, http.StatusInternalServerError, "internal_error", "the request could not be completed")
	}
}

func actorID(c *gin.Context) int64 {
	return c.GetInt64(middleware.ContextUserID)
}

func actorRole(c *gin.Context) string {
	return c.GetString(middleware.ContextRole)
}
