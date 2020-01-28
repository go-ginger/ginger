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

func Iterate(request m.IRequest, data interface{}, strict bool) (errors map[string]*me.ErrorItem) {
	s, ok := data.(reflect.Value)
	if !ok {
		s = reflect.ValueOf(data).Elem()
	}
	errors = make(map[string]*me.ErrorItem)
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
			validate, ok := ff.Tag.Lookup("validation")
			if ok {
				items := strings.Split(validate, ",")
				for _, item := range items {
					if validator, ok := Validators[item]; ok {
						errItem := validator.Handle(request, &ff, &f, nil)
						if errItem != nil {
							errors[ff.Name] = errItem
							break
						}
					}
				}
			}
		}
		break
	}
	return
}
