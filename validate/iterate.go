package validate

import (
	m "github.com/go-ginger/models"
	me "github.com/go-ginger/models/errors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"reflect"
	"strings"
)

func setNestedErrors(request m.IRequest, errors, nestedErrors map[string]*me.ErrorItem, field *reflect.StructField) {
	errors[field.Name] = &me.ErrorItem{
		Key: field.Name,
		Title: request.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    field.Name,
				Other: field.Name,
			},
		}),
		Errors: nestedErrors,
	}
}

func handleValidation(request m.IRequest, structField *reflect.StructField, field *reflect.Value,
	tagName string, errors map[string]*me.ErrorItem) (errorsResult map[string]*me.ErrorItem, handled bool) {
	errorsResult = errors
	if errorsResult == nil {
		errorsResult = make(map[string]*me.ErrorItem)
	}
	validate, ok := structField.Tag.Lookup(tagName)
	if ok {
		items := strings.Split(validate, ",")
		for _, item := range items {
			parts := strings.Split(item, "=")
			var tagValue *string
			tagName := parts[0]
			if len(parts) > 1 {
				tagValue = &parts[1]
			}
			if validator, ok := Validators[tagName]; ok {
				errItem := validator.Handle(request, structField, field, tagValue)
				if errItem != nil {
					errors[structField.Name] = errItem
					break
				}
			}
		}
	}
	return
}

func Iterate(request m.IRequest, data interface{}, strict bool) (errors map[string]*me.ErrorItem) {
	s, ok := data.(reflect.Value)
	if !ok {
		s = reflect.ValueOf(data).Elem()
	}
	sType := s.Type()
	switch s.Kind() {
	case reflect.Struct:
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			ff := sType.Field(i)
			switch f.Type().Kind() {
			case reflect.Ptr:
				if f.IsNil() {
					break
				}
				nestedErrors := Iterate(request, f.Elem(), strict)
				if nestedErrors != nil && len(nestedErrors) > 0 {
					setNestedErrors(request, errors, nestedErrors, &ff)
					return
				}
				break
			case reflect.Struct:
				nestedErrors := Iterate(request, f, strict)
				if nestedErrors != nil && len(nestedErrors) > 0 {
					setNestedErrors(request, errors, nestedErrors, &ff)
					return
				}
				break
			case reflect.Slice:
				for ind := 0; ind < f.Len(); ind++ {
					nestedErrors := Iterate(request, f.Index(ind), strict)
					if nestedErrors != nil && len(nestedErrors) > 0 {
						setNestedErrors(request, errors, nestedErrors, &ff)
						return
					}
				}
				break
			}
			if errors, ok = handleValidation(request, &ff, &f, "validation", errors); ok {
				break
			}
			ctx := request.GetContext()
			validationKey := ""
			switch ctx.Request.Method {
			case "POST":
				validationKey = "post_validation"
				break
			case "PUT":
				validationKey = "put_validation"
				break
			}
			if validationKey != "" {
				if errors, ok = handleValidation(request, &ff, &f, validationKey, errors); ok {
					break
				}
			}
		}
		break
	}
	return
}
