package prices

import (
	"context"
	"log"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/sptGabriel/investment-analyzer/extensions/pgtest"
	"github.com/sptGabriel/investment-analyzer/telemetry/logging"
)

var testCtx context.Context

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	ctx := context.Background()
	testCtx = logging.WithContext(ctx, zap.NewNop())

	teardown, err := pgtest.StartDockerContainer(pgtest.Config{})
	if err != nil {
		log.Panicf("starting docker container: %s", err)
	}

	defer teardown()

	return m.Run()
}
