package validate

import (
	"github.com/go-ginger/models"
	me "github.com/go-ginger/models/errors"
	"reflect"
)

type IValidator interface {
	Initialize()
	GetKey() *string
	Handle(request models.IRequest, field *reflect.StructField, value *reflect.Value, tagValue *string) (err *me.ErrorItem)
}
