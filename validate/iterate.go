package validate

import (
	"github.com/go-ginger/models"
	"reflect"
	"strings"
)

func addError(errors map[string]interface{}, name string, err error) {
	var fieldErrors []string
	if _, ok := errors[name]; !ok {
		fieldErrors = make([]string, 0)
		errors[name] = fieldErrors
	} else {
		fieldErrors = errors[name].([]string)
	}
	fieldErrors = append(fieldErrors, err.Error())
	errors[name] = fieldErrors
}

func Iterate(request models.IRequest, data interface{}, strict bool) (errors map[string]interface{}, err error) {
	s, ok := data.(reflect.Value)
	if !ok {
		s = reflect.ValueOf(data).Elem()
	}
	errors = make(map[string]interface{})
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
				nestedErrors, e := Iterate(request, f.Elem(), strict)
				if e != nil {
					err = e
					if nestedErrors != nil {
						addError(errors, ff.Name, err)
					}
					return
				}
				break
			case reflect.Struct:
				nestedErrors, e := Iterate(request, f, strict)
				if e != nil {
					err = e
					if nestedErrors != nil {
						addError(errors, ff.Name, err)
					}
					return
				}
				break
			case reflect.Slice:
				for ind := 0; ind < f.Len(); ind++ {
					nestedErrors, e := Iterate(request, f.Index(ind), strict)
					if e != nil {
						err = e
						if nestedErrors != nil {
							addError(errors, ff.Name, err)
						}
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
						err = validator.Handle(request, &ff, &f, nil)
						if err != nil {
							addError(errors, ff.Name, err)
							return
						}
					}
				}
			}
		}
		break
	}
	return
}
