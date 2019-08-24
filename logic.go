package ginger

import "github.com/kulichak/models"

type ILogic interface {
	Paginate(request models.IRequest)
}

type BaseLogic struct {
	ILogic
}

func (base *BaseLogic) Paginate(request models.IRequest) {

}
