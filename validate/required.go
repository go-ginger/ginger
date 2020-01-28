package validate

import (
	"github.com/go-ginger/models"
	"github.com/go-ginger/models/errors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"reflect"
)

type Required struct {
	IValidator
	Key *string
}

func (v *Required) Initialize() {
	if v.Key == nil {
		key := "required"
		v.Key = &key
	}
}

func (v *Required) GetKey() *string {
	return v.Key
}

func (v *Required) Handle(request models.IRequest, field *reflect.StructField, value *reflect.Value,
	tagValue *string) (err *errors.ErrorItem) {
	if field == nil || value == nil {
		return
	}
	if !value.IsValid() {
		return
	}
	kind := value.Kind()
	if kind == reflect.Interface || kind == reflect.Ptr {
		if value.IsNil() {
			err = &errors.ErrorItem{
				Key: field.Name,
				Title: request.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    field.Name,
						Other: field.Name,
					},
				}),
				Message: request.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "FieldRequired",
						Other: "Missing data for required field.",
					},
				}),
			}
		}
	}
	return
}
