package v1

import (
	"github.com/go-chi/chi/v5"
	swhttp "github.com/swaggo/http-swagger"

	swdoc "github.com/sptGabriel/investment-analyzer/docs/swagger" // Swagger content.
	"github.com/sptGabriel/investment-analyzer/extensions/gbhttp/rest"
)

// API
//
//	@Title			Investment Analyzer
//	@Description	Investment Analyzer REST API.
//	@Version		0.0.1
//	@License.name	Stone Co®
//	@Schemes		http
type API struct {
	ReportHandler reportHandler
}

func (a API) Routes(router chi.Router) {
	router.Get("/docs/v1/investment_analyzer/swagger/*", swhttp.Handler(
		swhttp.InstanceName(swdoc.SwaggerInfo.InstanceName()),
	))

	router.Post("/api/v1/investment_analyzer/portfolios/{portfolio-id}/reports", rest.Handle(a.ReportHandler.Handler))
}
