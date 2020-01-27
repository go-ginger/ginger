package ginger

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-ginger/ginger/validate"
	"github.com/go-ginger/models"
	gme "github.com/go-ginger/models/errors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (c *BaseController) handleRequestBody(ctx *gin.Context, request models.IRequest) (err error) {
	if c.DbHandler != nil {
		if ctx.Request.ContentLength > 0 {
			model := c.DbHandler.GetModelInstance().(models.IBaseModel)
			err = BindJSON(ctx, model)
			if err != nil {
				if c.ValidateRequestBody != nil && *c.ValidateRequestBody {
					err = errors.New(request.MustLocalize(&i18n.LocalizeConfig{
						DefaultMessage: &i18n.Message{
							ID:    "InvalidRequestInformation",
							Other: "Invalid request information. error: {{.Error}}",
						},
						TemplateData: map[string]string{
							"Error": err.Error(),
						},
					}))
					return
				}
			} else {
				validationErrors, e := validate.Iterate(request, model, c.StrictValidation)
				if e != nil {
					err = e
					if validationErrors != nil {
						if currentErr, ok := err.(gme.IError); ok {
							currentErr.SetErrors(validationErrors)
						}
					}
					return
				}
				request.SetBody(model)
			}
		}
	}
	return
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
	request := &models.Request{
		Context: ctx,
		Params:  &models.Params{},
		ID:      ctx.Params.ByName("id"),
		Filters: &filters,
		Sort:    &sort,
		Fields:  &fields,
	}
	for _, param := range ctx.Params {
		request.Params.Set(&models.Param{
			Key:   param.Key,
			Value: param.Value,
		})
	}
	if request.ID != "" {
		if request.Filters == nil || *request.Filters == nil {
			request.Filters = &models.Filters{}
		}
		(*request.Filters)["id"] = request.ID
	}
	acceptLanguage := ctx.GetHeader("Accept-Language")
	request.CurrentLanguage = &models.Language{
		AcceptLanguage: acceptLanguage,
		Localizer:      i18n.NewLocalizer(config.LanguageBundle, acceptLanguage),
	}
	sample := c.Controller.GetRequestSample()
	sample.SetBaseRequest(request)
	err := c.Controller.HandleRequestBody(ctx, sample)
	if err != nil {
		return sample, err
	}
	request = sample.GetBaseRequest()
	return sample, nil
}
