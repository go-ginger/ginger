package ginger

import (
	"github.com/gin-gonic/gin"
	"github.com/kulichak/logic"
	"github.com/kulichak/models"
)

type IController interface {
	Init()
	get(ctx *gin.Context)
	Post(request *models.Request)
	Get(request *models.Request)
	Put(request *models.Request)
	Delete(request *models.Request)
}

type BaseController struct {
	IController

	LogicHandler logic.IBaseLogicHandler
}

func (c *BaseController) Init() {
}

func (c *BaseController) HandleError(request *models.Request, result interface{}, err error) (handled bool) {
	if err != nil {
		request.Context.JSON(400, err)
		return true
	} else if result == nil {
		request.Context.JSON(404, result)
		return true
	}
	return false
}

func (c *BaseController) handleFilters(ctx *gin.Context) {
	ctx.Set("filters", GetQueryFilters(ctx))
}

func (c *BaseController) handlePagination(ctx *gin.Context) {
	ctx.Set("sort", GetSortFields(ctx))
	ctx.Set("page", ctx.Query("page"))
	ctx.Set("per_page", ctx.Query("per_page"))
}

func (c *BaseController) handleFields(ctx *gin.Context) {
	ctx.Set("fields", GetFetchFields(ctx, nil))
}

func (c *BaseController) get(ctx *gin.Context) {
	c.handleFields(ctx)
	c.handleFilters(ctx)
}
