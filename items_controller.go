package ginger

import (
	"github.com/gin-gonic/gin"
	"github.com/kulichak/models"
)

type IBaseItemsController interface {
	IController
}

type BaseItemsController struct {
	IBaseItemsController
	BaseController

	Controller IBaseItemsController
}

func (c *BaseItemsController) GetRoutes() []BaseControllerRoute {
	return c.Routes
}

func (c *BaseItemsController) RegisterRoutes(controller IBaseItemsController, path string, router *RouterGroup) {
	c.Controller = controller
	router.RegisterRoutes(controller, path, router.RouterGroup)
}

// POST
func (c *BaseItemsController) Post(request models.IRequest) {
	result, err := c.LogicHandler.DoCreate(request)
	if c.HandleError(request, result, err) {
		return
	}
	req := request.GetBaseRequest()
	req.Context.JSON(201, result)
}

func (c *BaseItemsController) post(ctx *gin.Context) {
	c.BaseController.post(ctx)
	req, err := c.NewRequest(ctx)
	if c.HandleErrorNoResult(req, err) {
		return
	}
	c.Controller.Post(req)
}

// GET
func (c *BaseItemsController) Get(request models.IRequest) {
	result, err := c.LogicHandler.DoPaginate(request)
	if c.HandleError(request, result, err) {
		return
	}
	req := request.GetBaseRequest()
	req.Context.JSON(200, result)
}

func (c *BaseItemsController) get(ctx *gin.Context) {
	c.handlePagination(ctx)
	c.BaseController.get(ctx)
	req, err := c.NewRequest(ctx)
	if c.HandleErrorNoResult(req, err) {
		return
	}
	c.Controller.Get(req)
}

// PUT
func (c *BaseItemsController) Put(request models.IRequest) {
	err := c.LogicHandler.DoUpdate(request)
	if c.HandleError(request, nil, err) {
		return
	}
	req := request.GetBaseRequest()
	req.Context.JSON(204, nil)
}

func (c *BaseItemsController) put(ctx *gin.Context) {
	c.BaseController.put(ctx)
	req, err := c.NewRequest(ctx)
	if c.HandleErrorNoResult(req, err) {
		return
	}
	c.Controller.Put(req)
}
