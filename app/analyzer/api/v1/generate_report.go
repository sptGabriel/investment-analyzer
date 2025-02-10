package v1

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/sptGabriel/investment-analyzer/app/analyzer/api/v1/scheme"
	"github.com/sptGabriel/investment-analyzer/domain/reports"
	"github.com/sptGabriel/investment-analyzer/extensions/gbhttp/rest"
	"github.com/sptGabriel/investment-analyzer/extensions/gblib"
)

const dateLayout = "2006-01-02 15:04:05"

type GenerateReportRequest struct {
	StartDate   string  `json:"start_date" validate:"required"`
	EndDate     string  `json:"end_date" validate:"required"`
	Interval    string  `json:"interval" validate:"required"`
	InitialCash float64 `json:"initial_cash" validate:"required"`
}

// Generate Report
//
//	@Summary		Generate a report
//	@Description	Generates a report based on specified criteria.  This endpoint allows users to retrieve data in a formatted report.  The report parameters (e.g., date range, internval) should be provided in the request body.
//	@Tags			Service
//	@Accept			json
//	@Produce		json
//	@Param			Body	body		GenerateReportRequest	true	"Body"
//	@Success		200		{object}	scheme.GenerateReportResponse
//	@Failure		400		{object}	rest.ErrorPayload
//	@Failure		500		{object}	rest.ErrorPayload
//	@Router			/api/v1/investment_analyzer/portfolios/{portfolio-id}/reports/ [POST]
func (h reportHandler) Handler(r *http.Request) rest.Response {
	reqBody, err := rest.ParseBody[GenerateReportRequest](r)
	if err != nil {
		return rest.BadRequest(err, rest.NewErrorPayload(rest.TypeBadRequest, "invalid body content"))
	}

	if err := h.validator.Struct(reqBody); err != nil {
		return rest.BadRequest(
			err, rest.NewErrorPayload(rest.TypeBadRequest, "invalid body content"))
	}

	dur, err := time.ParseDuration(reqBody.Interval)
	if err != nil {
		return rest.BadRequest(
			err, rest.NewErrorPayload(rest.TypeBadRequest, "invalid interval format"))
	}

	startDate, err := time.Parse(dateLayout, reqBody.StartDate)
	if err != nil {
		return rest.BadRequest(
			err, rest.NewErrorPayload(rest.TypeBadRequest, "invalid date format"))
	}

	endDate, err := time.Parse(dateLayout, reqBody.EndDate)
	if err != nil {
		return rest.BadRequest(
			err, rest.NewErrorPayload(rest.TypeBadRequest, "invalid date format"))
	}

	output, err := h.uc.Execute(r.Context(), reports.GenerateReportInput{
		PortfolioID: "408186c6-b76a-4ad6-8d4a-9ace3762b997", // default portfolio to challenge
		StartDate:   startDate,
		EndDate:     endDate,
		Interval:    dur,
	})
	if err != nil {
		return rest.InternalServerError(err)
	}

	return rest.OK(scheme.BuildGenerateReportResponse(output))
}

type reportHandler struct {
	uc        gblib.UseCase[reports.GenerateReportInput, reports.GenerateReportOutput]
	validator *validator.Validate
}

func NewReportHandler(
	uc gblib.UseCase[reports.GenerateReportInput, reports.GenerateReportOutput],
) reportHandler {
	return reportHandler{
		uc:        uc,
		validator: validator.New(),
	}
}
