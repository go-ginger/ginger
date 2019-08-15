package ginger

import (
	"github.com/gin-gonic/gin"
	"github.com/kulichak/logic"
)

type IController interface {
	Init()
	get(ctx *gin.Context)
}

type BaseController struct {
	IController

	LogicHandler logic.IBaseLogicHandler
}

func (c *BaseController) Init() {
}

func (c *BaseController) HandleError(request *Request, result interface{}, err error) (handled bool) {
	if err != nil {
		request.Context.JSON(400, err)
		return true
	} else if result == nil {
		request.Context.JSON(404, result)
		return true
	}
	return false
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


func (c *BaseItemsController) Init() {
}

func (c *BaseItemsController) RegisterRoutes(controller IBaseItemsController, path string, router *RouterGroup) {
	c.Controller = controller
	router.RegisterRoutes(controller, path, router.RouterGroup)
}

func (c *BaseItemsController) Get(request *Request) {
	result, err := c.LogicHandler.Paginate(&request.Request.IRequest)
	if c.HandleError(request, result, err) {
		return
	}
	request.Context.JSON(200, result)
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


func (c *BaseItemController) Init() {
}

func (c *BaseItemController) RegisterRoutes(controller IBaseItemController, path string, router *RouterGroup) {
	c.Controller = controller
	router.RegisterRoutes(controller, path, router.RouterGroup)
}

func (c *BaseItemController) get(ctx *gin.Context) {
	c.BaseController.get(ctx)
	c.Controller.Get(NewRequest(ctx))
}
