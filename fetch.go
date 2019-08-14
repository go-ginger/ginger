package ginger

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kulichak/ginger/helpers"
	"github.com/kulichak/models"
	"strings"
)

func GetQueryFilters(ctx *gin.Context) map[string]interface{} {
	filters := ctx.Query("filters")
	var filter map[string]interface{}
	if filters != "" {
		filter = map[string]interface{}{}
		json.Unmarshal([]byte(filters), filter)
	}
	return filter
}

func GetSortFields(ctx *gin.Context) []models.SortItem {
	sort := ctx.Query("sort")
	if sort == "" {
		sort = "-_score"
	}
	sorts := strings.Split(sort, ",")
	result := make([]models.SortItem, 0)
	for _, sort := range sorts {
		asc := true
		if strings.HasPrefix(sort, "-") {
			asc = false
			sort = sort[1:]
		}
		result = append(result, models.SortItem{
			Name:      sort,
			Ascending: asc,
		})
	}
	return result
}

func GetFetchFields(ctx *gin.Context, allowedFields []string) []string {
	fields := ctx.Query("fields")
	result := make([]string, 0)
	if fields != "" {
		for _, field := range strings.Split(fields, ",") {
			if allowedFields == nil || helpers.Contains(allowedFields, field) {
				result = append(result, field)
			}
		}
		return result
	}
	return nil
}
