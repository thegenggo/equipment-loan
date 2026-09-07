package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thegenggo/equipment-loan/api/internal/httperr"
	"github.com/thegenggo/equipment-loan/api/internal/service"
)

type createEquipmentRequest struct {
	Code     string `json:"code" binding:"required,max=50"`
	Name     string `json:"name" binding:"required,max=150"`
	Category string `json:"category" binding:"required,max=50"`
}

type updateEquipmentRequest struct {
	Code     string `json:"code" binding:"required,max=50"`
	Name     string `json:"name" binding:"required,max=150"`
	Category string `json:"category" binding:"required,max=50"`
	Status   string `json:"status" binding:"required,oneof=available repair"`
}

type EquipmentHandler struct {
	equipments *service.EquipmentService
}

func NewEquipmentHandler(equipments *service.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{equipments: equipments}
}

func (h *EquipmentHandler) List(c *gin.Context) {
	equipments, err := h.equipments.List(c.Request.Context(), c.Query("status"))
	if err != nil {
		if errors.Is(err, service.ErrInvalidStatus) {
			httperr.Write(c, http.StatusBadRequest, "invalid_status", "status must be available, borrowed or repair")
			return
		}
		log.Printf("list equipments: %v", err)
		httperr.Write(c, http.StatusInternalServerError, "internal_error", "could not load the catalogue")
		return
	}

	c.JSON(http.StatusOK, equipments)
}

func (h *EquipmentHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}

	equipment, err := h.equipments.Get(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err, "get equipment")
		return
	}

	c.JSON(http.StatusOK, equipment)
}

func (h *EquipmentHandler) Create(c *gin.Context) {
	var req createEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.Write(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	equipment, err := h.equipments.Create(c.Request.Context(), req.Code, req.Name, req.Category)
	if err != nil {
		h.writeError(c, err, "create equipment")
		return
	}

	c.JSON(http.StatusCreated, equipment)
}

func (h *EquipmentHandler) Update(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}

	var req updateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.Write(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	equipment, err := h.equipments.Update(c.Request.Context(), id, req.Code, req.Name, req.Category, req.Status)
	if err != nil {
		h.writeError(c, err, "update equipment")
		return
	}

	c.JSON(http.StatusOK, equipment)
}

func (h *EquipmentHandler) Delete(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}

	if err := h.equipments.Delete(c.Request.Context(), id); err != nil {
		h.writeError(c, err, "delete equipment")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *EquipmentHandler) writeError(c *gin.Context, err error, operation string) {
	switch {
	case errors.Is(err, service.ErrEquipmentNotFound):
		httperr.Write(c, http.StatusNotFound, "not_found", "equipment not found")
	case errors.Is(err, service.ErrEquipmentCodeTaken):
		httperr.Write(c, http.StatusConflict, "code_taken", "this equipment code is already is use")
	case errors.Is(err, service.ErrEquipmentBorrowed):
		httperr.Write(c, http.StatusConflict, "equipment_borrowed", "equipment on loan cannot be edited")
	case errors.Is(err, service.ErrEquipmentInUse):
		httperr.Write(c, http.StatusConflict, "equipment_in_use", "equipment with loan history cannot be deleted")
	default:
		log.Printf("%s: %v", operation, err)
		httperr.Write(c, http.StatusInternalServerError, "internal_error", "the request could not be completed")
	}
}

func pathID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		httperr.Write(c, http.StatusBadRequest, "invalid_id", "id must be a positive integer")
		return 0, false
	}

	return id, true
}
