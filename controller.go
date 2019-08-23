package ginger

import (
	"github.com/gin-gonic/gin"
	"github.com/kulichak/dl"
	"github.com/kulichak/logic"
	"github.com/kulichak/models"
)

type IController interface {
	GetRoutes() []BaseControllerRoute

	post(ctx *gin.Context)
	get(ctx *gin.Context)
	put(ctx *gin.Context)

	Post(request *models.Request)
	Get(request *models.Request)
	Put(request *models.Request)
	Delete(request *models.Request)
}

type BaseControllerRoute struct {
	Method   string
	Handlers []gin.HandlerFunc
}

type BaseController struct {
	IController

	Routes       []BaseControllerRoute
	LogicHandler logic.IBaseLogicHandler
}

func (c *BaseController) Init(logicHandler logic.IBaseLogicHandler, dbHandler dl.IBaseDbHandler) {
	c.LogicHandler = logicHandler
	c.LogicHandler.Init(dbHandler)
}

func (c *BaseController) AddRoute(method string, handlers ...gin.HandlerFunc) {
	c.Routes = append(c.Routes, BaseControllerRoute{
		Method:   method,
		Handlers: handlers,
	})
}

func (c *BaseController) GetRoutes() []BaseControllerRoute {
	return c.Routes
}

func (c *BaseController) handleError(err error) (*int, error) {
	if err != nil {
		status := 400
		message := err.Error()
		if e, ok := err.(models.Error); ok {
			status = e.Status
		}
		return &status, &models.Error{
			Status:  status,
			Message: message,
		}
	}
	return nil, nil
}

func (c *BaseController) HandleErrorNoResult(request *models.Request, err error) (handled bool) {
	if err != nil {
		status, e := c.handleError(err)
		if status != nil && e != nil {
			request.Context.JSON(*status, models.Error{
				Message: e.Error(),
			})
			return true
		}
	}
	return false
}

func (c *BaseController) HandleError(request *models.Request, result interface{}, err error) (handled bool) {
	if err != nil {
		status, e := c.handleError(err)
		if status != nil && e != nil {
			request.Context.JSON(*status, models.Error{
				Message: e.Error(),
			})
			return true
		}
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
	if _, ok := ctx.GetQuery("sort"); ok {
		ctx.Set("sort", GetSortFields(ctx))
	}
	queries := []string{"page", "per_page"}
	for _, query := range queries {
		if q, ok := ctx.GetQuery(query); ok {
			ctx.Set(query, q)
		}
	}
}

func (c *BaseController) handleFields(ctx *gin.Context) {
	ctx.Set("fields", GetFetchFields(ctx, nil))
}

func (c *BaseController) post(ctx *gin.Context) {
}

func (c *BaseController) get(ctx *gin.Context) {
	c.handleFields(ctx)
	c.handleFilters(ctx)
}

func (c *BaseController) put(ctx *gin.Context) {
}
