package gin_extended

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	RouterGroup

	Engine  *gin.Engine
	Address []string
}

func NewRouter() *Router  {
	router := new(Router)
	router.Engine = gin.Default()
	router.RouterGroup.engine = router.Engine
	return router
}

func(router Router) Run() (err error) {
	return router.Engine.Run(router.Address...)
}
