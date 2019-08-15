package ginger

import "github.com/gin-gonic/gin"

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

