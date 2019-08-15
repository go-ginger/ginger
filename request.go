package ginger

import (
	"github.com/gin-gonic/gin"
	"github.com/kulichak/ginger/helpers"
	"github.com/kulichak/models"
	"strconv"
)

type Request struct {
	models.Request
}

var methodsWithBody []string

func init() {
	methodsWithBody = []string{"POST", "PUT"}
}

func NewRequest(ctx *gin.Context) *Request {
	filtersFace, exists := ctx.Get("filters")
	var filters models.Filters
	if exists {
		filters = filtersFace.(map[string]interface{})
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
	perPageFace, exists := ctx.Get("per_page")
	var perPage uint64
	if exists {
		perPage, _ = strconv.ParseUint(perPageFace.(string), 10, 32)
	}
	request := &Request{
		Request: models.Request{
			Context: ctx,
			Params:  &ctx.Params,
			ID:      ctx.Params.ByName("id"),
			Filters: &filters,
			Sort:    &sort,
			Fields:  &fields,
			Page:    &page,
			PerPage: &perPage,
		},
	}
	if helpers.Contains(methodsWithBody, ctx.Request.Method) {
		ctx.BindJSON(request.Body)
	}
	return request
}
