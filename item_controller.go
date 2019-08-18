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
	c.Controller.Post(c.NewRequest(ctx))
}

// GET
func (c *BaseItemController) Get(request *models.Request) {
	result, err := c.LogicHandler.DoGet(request)
	if c.HandleError(request, result, err) {
		return
	}
	request.Context.JSON(200, result)
}

func (c *BaseItemController) get(ctx *gin.Context) {
	c.BaseController.get(ctx)
	req := c.NewRequest(ctx)
	req.Filters = &models.Filters{
		"id": req.ID,
	}
	c.Controller.Get(req)
}

// PUT
func (c *BaseItemController) Put(request *models.Request) {
	err := c.LogicHandler.DoUpdate(request)
	if c.HandleErrorNoResult(request, err) {
		return
	}
	request.Context.JSON(204, nil)
}

func (c *BaseItemController) put(ctx *gin.Context) {
	c.BaseController.put(ctx)
	c.Controller.Put(c.NewRequest(ctx))
}
