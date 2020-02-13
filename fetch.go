package ginger

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-ginger/helpers"
	"github.com/go-ginger/models"
	gm "github.com/go-ginger/models"
	"reflect"
	"strings"
)

func GetQueryFilters(ctx *gin.Context) map[string]interface{} {
	filters := ctx.Query("filters")
	var filter map[string]interface{}
	if filters != "" {
		filter = map[string]interface{}{}
		json.Unmarshal([]byte(filters), &filter)
	}
	return filter
}

func GetSortFields(ctx *gin.Context) []models.SortItem {
	sort := ctx.Query("sort")
	result := make([]models.SortItem, 0)
	if sort == "" {
		return result
	}
	sorts := strings.Split(sort, ",")
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

func (c *BaseController) clear(value *reflect.Value, _type *reflect.Type) {
	if value == nil {
		return
	}
	if value.IsValid() {
		if value.CanSet() {
			value.Set(reflect.Zero(*_type))
		}
	}
}

type iBeforeDump interface {
	BeforeDump(request gm.IRequest, data interface{})
}

func (c *BaseController) BeforeDump(request gm.IRequest, data interface{}) {
	s, ok := data.(reflect.Value)
	if !ok {
		s = reflect.ValueOf(data)
	}
	kind := s.Kind()
	if kind == reflect.Ptr {
		s = s.Elem()
		kind = s.Kind()
	}
	sType := s.Type()
	switch kind {
	case reflect.Struct:
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			ff := sType.Field(i)
			if helpers.IsEmptyValue(f) {
				continue
			}
			switch f.Type().Kind() {
			case reflect.Ptr:
				if f.IsNil() {
					break
				}
				c.BeforeDump(request, f.Elem())
				break
			case reflect.Struct:
				c.BeforeDump(request, f)
				break
			case reflect.Slice:
				for ind := 0; ind < f.Len(); ind++ {
					c.BeforeDump(request, f.Index(ind))
				}
				break
			}
			if cs, ok := ff.Tag.Lookup("c"); ok {
				continueCheck := true
				csParts := strings.Split(cs, ",")
				for _, csPart := range csParts {
					if csPart == "load_only" {
						continueCheck = false
						break
					}
				}
				if !continueCheck {
					continue
				}
			}
			tag, ok := ff.Tag.Lookup("read_roles")
			if ok {
				canRead := false
				auth := request.GetAuth()
				if auth != nil {
					tagParts := strings.Split(tag, ",")
					for _, role := range tagParts {
						if auth.HasRole(role) || (role == "id" &&
							auth.GetCurrentAccountId(request) == request.GetIDString()) {
							canRead = true
							break
						}
					}
				}
				if !canRead {
					c.clear(&f, &ff.Type)
				}
			}
		}
		if s.CanAddr() {
			addr := s.Addr()
			if addr.IsValid() && addr.CanInterface() {
				mv := addr.Interface()
				if baseModel, ok := mv.(gm.IBaseModel); ok {
					baseModel.Populate(request)
				}
				if cls, ok := mv.(iBeforeDump); ok {
					cls.BeforeDump(request, mv)
				}
			}
		}
		break
	}
}
