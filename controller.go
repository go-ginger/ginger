package gin_extended

import (
	"github.com/gin-gonic/gin"
)

type ControllerFace interface {
	get(ctx *gin.Context)
}

type BaseController struct {
	ControllerFace
}

func (c *BaseController) handleFilters(ctx *gin.Context) {
	ctx.Set("filters", GetQueryFilters(ctx))
}

func (c *BaseController) handlePagination(ctx *gin.Context) {
	ctx.Set("sort", GetSortFields(ctx))
	ctx.Set("page", ctx.Query("page"))
	ctx.Set("per_page", ctx.Query("per_page"))
}

func (c *BaseController) handleFields(ctx *gin.Context) {
	ctx.Set("fields", GetFetchFields(ctx, nil))
}

func (c *BaseController) get(ctx *gin.Context) {
	c.handleFields(ctx)
	c.handleFilters(ctx)
}

// items
type BaseItemsControllerFace interface {
	ControllerFace

	Get(request *Request)
}

// items controller
type BaseItemsController struct {
	BaseItemsControllerFace
	BaseController

	Controller BaseItemsControllerFace
}

func (c *BaseItemsController) RegisterRoutes(controller BaseItemsControllerFace, path string, router *RouterGroup) {
	c.Controller = controller
	router.RegisterRoutes(controller, path, router.RouterGroup)
}

func (c *BaseItemsController) get(ctx *gin.Context) {
	c.handlePagination(ctx)
	c.BaseController.get(ctx)
	c.Controller.Get(NewRequest(ctx))
}

// item
type BaseItemControllerFace interface {
	ControllerFace

	Get(request *Request)
}

// item controller
type BaseItemController struct {
	BaseItemControllerFace
	BaseController

	Controller BaseItemControllerFace
}

func (c *BaseItemController) RegisterRoutes(controller BaseItemControllerFace, path string, router *RouterGroup) {
	c.Controller = controller
	router.RegisterRoutes(controller, path, router.RouterGroup)
}

func (c *BaseItemController) get(ctx *gin.Context) {
	c.BaseController.get(ctx)
	c.Controller.Get(NewRequest(ctx))
}
