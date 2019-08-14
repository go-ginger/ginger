package gin_extended

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Request struct {
	Context *gin.Context
	Params  *gin.Params
	Fields  *[]string
	Filters *map[string]interface{}
	Sort    *[]SortItem
	Page    *int
	PerPage *int
}

func NewRequest(ctx *gin.Context) *Request {
	filtersFace, exists := ctx.Get("filters")
	var filters map[string]interface{}
	if exists {
		filters = filtersFace.(map[string]interface{})
	}
	sortFace, exists := ctx.Get("sort")
	var sort []SortItem
	if exists {
		sort = sortFace.([]SortItem)
	}
	fieldsFace, exists := ctx.Get("fields")
	var fields []string
	if exists {
		fields = fieldsFace.([]string)
	}
	pageFace, exists := ctx.Get("page")
	var page int
	if exists {
		page, _ = strconv.Atoi(pageFace.(string))
	}
	perPageFace, exists := ctx.Get("per_page")
	var perPage int
	if exists {
		perPage, _ = strconv.Atoi(perPageFace.(string))
	}
	return &Request{
		Context: ctx,
		Params:  &ctx.Params,
		Filters: &filters,
		Sort:    &sort,
		Fields:  &fields,
		Page:    &page,
		PerPage: &perPage,
	}
}
