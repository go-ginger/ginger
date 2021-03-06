package ginger

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ginger/helpers"
	"strings"
)

type RouterGroup struct {
	*gin.RouterGroup

	engine         *gin.Engine
	beforeRequests []HandlerFunc
}

func (group *RouterGroup) Group(relativePath string) *RouterGroup {
	var g *gin.RouterGroup
	if group.RouterGroup != nil {
		g = group.RouterGroup.Group(relativePath)
	} else {
		g = group.engine.Group(relativePath)
	}
	return &RouterGroup{
		RouterGroup: g,
		engine:      group.engine,
	}
}

func (group *RouterGroup) Any(handler HandlerFunc) (result interface{}) {
	if group.beforeRequests == nil {
		group.beforeRequests = make([]HandlerFunc, 0)
	}
	group.beforeRequests = append(group.beforeRequests, handler)
	return
}

func (group *RouterGroup) RegisterRoutes(controller IController, path string, router *gin.RouterGroup) {
	routes := controller.GetRoutes()
	routesMap := map[string]HandlerFunc{
		"any":    controller.any,
		"get":    controller.get,
		"post":   controller.post,
		"put":    controller.put,
		"delete": controller.delete,
	}
	methodHandlerNameMap := map[string]string{
		"any":    "Any",
		"post":   "Post",
		"get":    "Get",
		"put":    "Put",
		"delete": "Delete",
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
		method := strings.ToLower(route.Method)
		if handler, ok := routesMap[method]; ok {
			methodName := methodHandlerNameMap[method]
			f := helpers.ReflectMethod(controller, methodName)
			if f != nil {
				var handlers []gin.HandlerFunc
				if config.CorsEnabled {
					handlers = append(handlers, CORS)
				}
				handlers = append(handlers, baseHandler)
				if group.beforeRequests != nil {
					for _, handler := range group.beforeRequests {
						handlers = append(handlers, controller.GetHandler(group, RouteHandler{
							Handler:  handler,
							CallBack: nil,
						}))
					}
				}
				for _, handler := range route.Handlers {
					handlers = append(handlers, controller.GetHandler(group, handler))
				}
				handlers = append(handlers, controller.GetHandler(group, RouteHandler{
					Handler:  handler,
					CallBack: nil,
				}))
				switch method {
				case "any":
					router.Any(path, handlers...)
					break
				case "post":
					router.POST(path, handlers...)
					break
				case "get":
					router.GET(path, handlers...)
					break
				case "put":
					router.PUT(path, handlers...)
					break
				case "delete":
					router.DELETE(path, handlers...)
					break
				}
			}
		}
	}
}

func (group *RouterGroup) RegisterRoute(controller IController, router *gin.RouterGroup, path, method string,
	customHandlers ...HandlerFunc) {
	routes := controller.GetRoutes()
	if config.CorsEnabled {
		router.OPTIONS(path, CORS)
	}
	baseHandler := controller.GetHandler(group, RouteHandler{
		Type:     -1,
		Handler:  nil,
		CallBack: nil,
	})
	for _, route := range routes {
		var handlers []gin.HandlerFunc
		if config.CorsEnabled {
			handlers = append(handlers, CORS)
		}
		handlers = append(handlers, baseHandler)
		if group.beforeRequests != nil {
			for _, handler := range group.beforeRequests {
				handlers = append(handlers, controller.GetHandler(group, RouteHandler{
					Handler:  handler,
					CallBack: nil,
				}))
			}
		}
		for _, handler := range route.Handlers {
			handlers = append(handlers, controller.GetHandler(group, handler))
		}
		for _, h := range customHandlers {
			handlers = append(handlers, controller.GetHandler(group, RouteHandler{
				Handler:  h,
				CallBack: nil,
			}))
		}
		method = strings.ToUpper(method)
		router.Handle(method, path, handlers...)
	}
}
