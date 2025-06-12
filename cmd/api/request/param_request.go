package request

type IdParamRequest struct {
	Id uint `param:"id" binding:"required"`
}
