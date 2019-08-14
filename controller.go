package ginger

import (
	"github.com/gin-gonic/gin"
)

type IController interface {
	get(ctx *gin.Context)
}

type BaseController struct {
	IController
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
type IBaseItemsController interface {
	IController

	Get(request *Request)
}

// items controller
type BaseItemsController struct {
	IBaseItemsController
	BaseController

	Controller IBaseItemsController
}

func (c *BaseItemsController) RegisterRoutes(controller IBaseItemsController, path string, router *RouterGroup) {
	c.Controller = controller
	router.RegisterRoutes(controller, path, router.RouterGroup)
}

func (c *BaseItemsController) get(ctx *gin.Context) {
	c.handlePagination(ctx)
	c.BaseController.get(ctx)
	c.Controller.Get(NewRequest(ctx))
}

// item
type IBaseItemController interface {
	IController

	Get(request *Request)
}

// item controller
type BaseItemController struct {
	IBaseItemController
	BaseController

	Controller IBaseItemController
}

func (c *BaseItemController) RegisterRoutes(controller IBaseItemController, path string, router *RouterGroup) {
	c.Controller = controller
	router.RegisterRoutes(controller, path, router.RouterGroup)
}

func (c *BaseItemController) get(ctx *gin.Context) {
	c.BaseController.get(ctx)
	c.Controller.Get(NewRequest(ctx))
}
