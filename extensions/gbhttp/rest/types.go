package rest

type Type string

const (
	TypeServerError         Type = "srn:error:server_error"
	TypeBadRequest          Type = "srn:error:invalid_params"
	TypeNotFound            Type = "srn:error:resource_not_found"
	TypeConflict            Type = "srn:error:conflict"
	TypeForbidden           Type = "srn:error:forbidden"
	TypeUnauthorized        Type = "srn:error:unauthorized"
	TypeNotImplemented      Type = "srn:error:not_implemented"
	TypeUnprocessableEntity Type = "srn:error:unprocessable_entity"
)
