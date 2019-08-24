package ginger

import (
	"github.com/gin-gonic/gin"
	"github.com/kulichak/models"
)

// item
type IBaseItemController interface {
	IController
}

// item controller
type BaseItemController struct {
	IBaseItemController
	BaseController

	Controller IBaseItemController
}

func (c *BaseItemController) GetRoutes() []BaseControllerRoute {
	return c.Routes
}

func (c *BaseItemController) RegisterRoutes(controller IBaseItemController, path string, router *RouterGroup) {
	c.Controller = controller
	router.RegisterRoutes(controller, path, router.RouterGroup)
}

// POST
func (c *BaseItemController) post(ctx *gin.Context) {
	c.BaseController.post(ctx)
	req, err := c.NewRequest(ctx)
	if c.HandleErrorNoResult(req, err) {
		return
	}
	c.Controller.Post(req)
}

// GET
func (c *BaseItemController) Get(request models.IRequest) {
	result, err := c.LogicHandler.DoGet(request)
	if c.HandleError(request, result, err) {
		return
	}
	req := request.GetBaseRequest()
	req.Context.JSON(200, result)
}

func (c *BaseItemController) get(ctx *gin.Context) {
	c.BaseController.get(ctx)
	req, err := c.NewRequest(ctx)
	if c.HandleErrorNoResult(req, err) {
		return
	}
	c.Controller.Get(req)
}

// PUT
func (c *BaseItemController) Put(request models.IRequest) {
	err := c.LogicHandler.DoUpdate(request)
	if c.HandleErrorNoResult(request, err) {
		return
	}
	req := request.GetBaseRequest()
	req.Context.JSON(204, nil)
}

func (c *BaseItemController) put(ctx *gin.Context) {
	c.BaseController.put(ctx)
	req, err := c.NewRequest(ctx)
	if c.HandleErrorNoResult(req, err) {
		return
	}
	c.Controller.Put(req)
}
