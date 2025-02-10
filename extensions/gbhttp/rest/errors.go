package rest

type errorPayloadDetail struct {
	Reason   string      `json:"reason" example:""`
	Value    interface{} `json:"value,omitempty" example:""`
	Property string      `json:"property,omitempty" example:""`
}

type ErrorPayload struct {
	Type    string               `json:"type" example:"srn:error:some_error"`
	Title   string               `json:"title,omitempty" example:"Message for some error"`
	Details []errorPayloadDetail `json:"details,omitempty" swaggertype:"object"`
}

func (e ErrorPayload) Detail(reason, value, property string) ErrorPayload {
	e.Details = append(e.Details, errorPayloadDetail{Reason: reason, Value: value, Property: property})
	return e
}

func NewErrorPayload(_type Type, title string) ErrorPayload {
	return ErrorPayload{
		Type:  string(_type),
		Title: title,
	}
}
