package ginger

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ginger/helpers"
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
	routesMap := map[string]HandlerFunc{
		"Get":  controller.get,
		"Post": controller.post,
		"Put":  controller.put,
	}
	if config.CorsEnabled {
		router.OPTIONS(path, CORS)
	}
	baseHandler := controller.GetHandler(group, RouteHandler{
		Type:     -1,
		Handler:  nil,
		CallBack: nil,
	})
	for _, route := range routes {
		if handler, ok := routesMap[route.Method]; ok {
			f := helpers.ReflectMethod(controller, route.Method)
			if f != nil {
				var handlers []gin.HandlerFunc
				if config.CorsEnabled {
					handlers = append(handlers, CORS)
				}
				handlers = append(handlers, baseHandler)
				for _, handler := range route.Handlers {
					handlers = append(handlers, controller.GetHandler(group, handler))
				}
				handlers = append(handlers, controller.GetHandler(group, RouteHandler{
					Handler:  handler,
					CallBack: nil,
				}))
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
