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

func (c *BaseItemController) GetHandler(group *RouterGroup, routeHandler RouteHandler) gin.HandlerFunc {
	return c.BaseController.GetHandler(group, routeHandler)
}

func (c *BaseItemController) GetRoutes() []BaseControllerRoute {
	return c.Routes
}

func (c *BaseItemController) RegisterRoutes(controller IBaseItemController, path string, router *RouterGroup) {
	c.Controller = controller
	router.RegisterRoutes(controller, path, router.RouterGroup)
}

// POST
func (c *BaseItemController) post(request models.IRequest) (result interface{}) {
	c.BaseController.post(request)
	result = c.Controller.Post(request)
	return
}

// GET
func (c *BaseItemController) Get(request models.IRequest) (result interface{}) {
	result, err := c.LogicHandler.DoGet(request)
	if c.HandleError(request, result, err) {
		return
	}
	req := request.GetBaseRequest()
	req.Context.JSON(200, result)
	return
}

func (c *BaseItemController) get(request models.IRequest) (result interface{}) {
	c.BaseController.get(request)
	result = c.Controller.Get(request)
	return
}

// PUT
func (c *BaseItemController) Put(request models.IRequest) (result interface{}) {
	err := c.LogicHandler.DoUpdate(request)
	if c.HandleErrorNoResult(request, err) {
		return
	}
	req := request.GetBaseRequest()
	req.Context.JSON(204, nil)
	return
}

func (c *BaseItemController) put(request models.IRequest) (result interface{}) {
	c.BaseController.put(request)
	result = c.Controller.Put(request)
	return
}
