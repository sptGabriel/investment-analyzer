package interceptors

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/extensions/gblib"
	"github.com/sptGabriel/investment-analyzer/extensions/utils"
)

type auditLogsRepository interface {
	Save(context.Context, entities.AuditLog) error
}

func AuditInterceptor(
	repository auditLogsRepository,
) gblib.Interceptor {
	return func(ctx context.Context, input interface{}, next gblib.InterceptorFunc) (interface{}, error) {
		startTime := time.Now()

		requestData, has := utils.RequestDataFromContext(ctx)
		if !has {
			requestData = utils.RequestData{
				IP:   "unk",
				Port: "unk",
			}
		}

		params, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("%w:on marshalling input", err)
		}

		auditLog := entities.AuditLog{
			Timestamp: startTime,
			IP:        requestData.IP,
			Port:      requestData.Port,
			Params:    string(params),
			Result:    "success",
		}

		result, err := next(ctx, input)
		if err != nil {
			auditLog.Result = "error"
			auditLog.ErrorMessage = err.Error()
		}

		if err := repository.Save(ctx, auditLog); err != nil {
			return nil, err
		}

		return result, err
	}
}
