package ginger

import (
	"github.com/gin-gonic/gin"
	"github.com/kulichak/ginger/helpers"
)

type RouterGroup struct {
	*gin.RouterGroup

	engine *gin.Engine
}

func (group *RouterGroup) Group(relativePath string) *RouterGroup {
	return &RouterGroup{
		RouterGroup: group.engine.Group(relativePath),
	}
}

func (group *RouterGroup) RegisterRoutes(controller IController, path string, router *gin.RouterGroup) {
	routes := map[string]gin.HandlerFunc{
		"Get": controller.get,
	}
	for method, handler := range routes {
		f := helpers.ReflectMethod(controller, method)
		if f != nil {
			switch method {
			case "Get":
				router.GET(path, handler)
				break
			}
		}
	}
}
