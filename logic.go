package ginger

type ILogic interface {
	Paginate(request *Request)
}

type BaseLogic struct {
	ILogic
}

func (base *BaseLogic) Paginate(request *Request) {

}
