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

func (c *BaseItemsController) RegisterRoutes(controller IBaseItemsController, path string, router *RouterGroup) {
	c.Controller = controller
	router.RegisterRoutes(controller, path, router.RouterGroup)
}

func (c *BaseItemsController) Post(request *models.Request) {
	result, err := c.LogicHandler.DoCreate(request)
	if c.HandleError(request, result, err) {
		return
	}
	request.Context.JSON(201, result)
}

func (c *BaseItemsController) post(ctx *gin.Context) {
	c.BaseController.post(ctx)
	c.Controller.Post(c.NewRequest(ctx))
}

func (c *BaseItemsController) Get(request *models.Request) {
	result, err := c.LogicHandler.DoPaginate(request)
	if c.HandleError(request, result, err) {
		return
	}
	request.Context.JSON(200, result)
}

func (c *BaseItemsController) get(ctx *gin.Context) {
	c.handlePagination(ctx)
	c.BaseController.get(ctx)
	c.Controller.Get(c.NewRequest(ctx))
}
