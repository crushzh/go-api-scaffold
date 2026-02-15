package handler

import (
	"strconv"

	"go-api-scaffold/internal/model"
	"go-api-scaffold/internal/service"
	"go-api-scaffold/pkg/response"

	"github.com/gin-gonic/gin"
)

// ExampleHandler handles Example CRUD endpoints
type ExampleHandler struct {
	svc *service.ExampleService
}

func NewExampleHandler(svc *service.ExampleService) *ExampleHandler {
	return &ExampleHandler{svc: svc}
}

// List returns a paginated list of examples
// @Summary  List examples
// @Tags     Example
// @Security Bearer
// @Param    page      query int    false "Page number"   default(1)
// @Param    page_size query int    false "Page size"     default(10)
// @Param    keyword   query string false "Search keyword"
// @Param    status    query string false "Status filter" Enums(active, inactive)
// @Success  200 {object} response.Response{data=response.PageData}
// @Router   /examples [get]
func (h *ExampleHandler) List(c *gin.Context) {
	var req model.QueryExampleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ParamError(c, "invalid parameters")
		return
	}

	items, total, err := h.svc.List(&req)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}

	response.SuccessPage(c, items, total, req.Page, req.PageSize)
}

// Create creates a new example
// @Summary  Create example
// @Tags     Example
// @Security Bearer
// @Accept   json
// @Produce  json
// @Param    body body model.CreateExampleRequest true "Create parameters"
// @Success  200  {object} response.Response{data=model.Example}
// @Router   /examples [post]
func (h *ExampleHandler) Create(c *gin.Context) {
	var req model.CreateExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "invalid parameters: "+err.Error())
		return
	}

	item, err := h.svc.Create(&req)
	if err != nil {
		response.ServerError(c, "create failed: "+err.Error())
		return
	}

	response.Success(c, item)
}

// Get returns an example by ID
// @Summary  Get example by ID
// @Tags     Example
// @Security Bearer
// @Param    id path int true "ID"
// @Success  200 {object} response.Response{data=model.Example}
// @Router   /examples/{id} [get]
func (h *ExampleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid ID")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "record not found")
		return
	}

	response.Success(c, item)
}

// Update updates an example
// @Summary  Update example
// @Tags     Example
// @Security Bearer
// @Accept   json
// @Param    id   path int                        true "ID"
// @Param    body body model.UpdateExampleRequest  true "Update parameters"
// @Success  200  {object} response.Response{data=model.Example}
// @Router   /examples/{id} [put]
func (h *ExampleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid ID")
		return
	}

	var req model.UpdateExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "invalid parameters: "+err.Error())
		return
	}

	item, err := h.svc.Update(uint(id), &req)
	if err != nil {
		response.ServerError(c, "update failed: "+err.Error())
		return
	}

	response.Success(c, item)
}

// Delete removes an example
// @Summary  Delete example
// @Tags     Example
// @Security Bearer
// @Param    id path int true "ID"
// @Success  200 {object} response.Response
// @Router   /examples/{id} [delete]
func (h *ExampleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid ID")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, "delete failed: "+err.Error())
		return
	}

	response.OK(c)
}
