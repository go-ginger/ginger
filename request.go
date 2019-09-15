package ginger

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kulichak/helpers"
	"github.com/kulichak/models"
	"strconv"
)

var methodsWithBody []string

func init() {
	methodsWithBody = []string{"POST", "PUT"}
}

func (c *BaseController) NewRequest(ctx *gin.Context) (models.IRequest, error) {
	filtersFace, exists := ctx.Get("filters")
	var filters models.Filters
	if exists {
		filters = filtersFace.(map[string]interface{})
	} else {
		filters = models.Filters{}
	}
	sortFace, exists := ctx.Get("sort")
	var sort []models.SortItem
	if exists {
		sort = sortFace.([]models.SortItem)
	}
	fieldsFace, exists := ctx.Get("fields")
	var fields models.Fields
	if exists {
		fields = fieldsFace.([]string)
	}
	pageFace, exists := ctx.Get("page")
	var page uint64
	if exists {
		page, _ = strconv.ParseUint(pageFace.(string), 10, 32)
	}
	if page <= 0 {
		page = 1
	}
	perPageFace, exists := ctx.Get("per_page")
	var perPage uint64
	if exists {
		perPage, _ = strconv.ParseUint(perPageFace.(string), 10, 32)
	}
	if perPage <= 0 {
		perPage = 30
	}
	request := &models.Request{
		Context: ctx,
		Params:  &ctx.Params,
		ID:      ctx.Params.ByName("id"),
		Filters: &filters,
		Sort:    &sort,
		Fields:  &fields,
		Page:    page,
		PerPage: perPage,
	}
	if request.ID != "" {
		if request.Filters == nil || *request.Filters == nil {
			request.Filters = &models.Filters{}
		}
		(*request.Filters)["id"] = request.ID
	}
	sample := c.Controller.GetRequestSample()
	sample.SetBaseRequest(request)
	request = sample.GetBaseRequest()
	if helpers.Contains(methodsWithBody, ctx.Request.Method) {
		if c.LogicHandler != nil {
			c.LogicHandler.Model(sample)
			sample.SetBody(request.Model)
			err := BindJSON(ctx, request.Body)
			if err != nil {
				return request, errors.New("Invalid request information. error: " + err.Error())
			}
			c.LogicHandler.Model(sample)
		}
	}
	return sample, nil
}
