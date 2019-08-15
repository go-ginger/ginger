package ginger

import "github.com/kulichak/models"

type ILogic interface {
	Paginate(request *models.Request)
}

type BaseLogic struct {
	ILogic
}

func (base *BaseLogic) Paginate(request *models.Request) {

}
