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

func (c *BaseItemsController) GetHandler(group *RouterGroup, routeHandler RouteHandler) gin.HandlerFunc {
	return c.BaseController.GetHandler(group, routeHandler)
}

func (c *BaseItemsController) GetRoutes() []BaseControllerRoute {
	return c.Routes
}

func (c *BaseItemsController) RegisterRoutes(controller IBaseItemsController, path string, router *RouterGroup) {
	c.Controller = controller
	router.RegisterRoutes(controller, path, router.RouterGroup)
}

// POST
func (c *BaseItemsController) Post(request models.IRequest) (result interface{}) {
	result, err := c.LogicHandler.DoCreate(request)
	if c.HandleError(request, result, err) {
		return
	}
	req := request.GetBaseRequest()
	req.Context.JSON(201, result)
	return
}

func (c *BaseItemsController) post(request models.IRequest) (result interface{}) {
	c.BaseController.post(request)
	result = c.Controller.Post(request)
	return
}

// GET
func (c *BaseItemsController) Get(request models.IRequest) (result interface{}) {
	result, err := c.LogicHandler.DoPaginate(request)
	if c.HandleError(request, result, err) {
		return
	}
	req := request.GetBaseRequest()
	req.Context.JSON(200, result)
	return
}

func (c *BaseItemsController) get(request models.IRequest) (result interface{}) {
	c.handlePagination(request)
	c.BaseController.get(request)
	c.Controller.Get(request)
	return
}

// PUT
func (c *BaseItemsController) Put(request models.IRequest) (result interface{}) {
	err := c.LogicHandler.DoUpdate(request)
	if c.HandleError(request, nil, err) {
		return
	}
	req := request.GetBaseRequest()
	req.Context.JSON(204, nil)
	return
}

func (c *BaseItemsController) put(request models.IRequest) (result interface{}) {
	c.BaseController.put(request)
	result = c.Controller.Put(request)
	return
}
