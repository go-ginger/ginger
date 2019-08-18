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
	routes := controller.GetRoutes()
	routesMap := map[string]gin.HandlerFunc{
		"Get":  controller.get,
		"Post": controller.post,
		"Put":  controller.put,
	}
	for _, route := range routes {
		if handler, ok := routesMap[route.Method]; ok {
			f := helpers.ReflectMethod(controller, route.Method)
			if f != nil {
				handlers := []gin.HandlerFunc{
					handler,
				}
				handlers = append(handlers, route.Handlers...)
				switch route.Method {
				case "Get":
					router.GET(path, handlers...)
					break
				case "Post":
					router.POST(path, handlers...)
					break
				case "Put":
					router.PUT(path, handlers...)
					break
				}
			}
		}
	}
}
