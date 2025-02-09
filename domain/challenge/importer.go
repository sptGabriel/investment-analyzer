package challenge

import (
	"context"
	"fmt"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/ports"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
	"github.com/sptGabriel/investment-analyzer/extensions/gbdb"
	"github.com/sptGabriel/investment-analyzer/extensions/gblib"
)

type settingsRepository interface {
	IsCsvImported(context.Context) (bool, error)
	SetCsvImported(context.Context) error
}

type assetsRepository interface {
	Upsert(context.Context, entities.Asset) error
}

type pricesRepository interface {
	Delete(context.Context) error
	SavePricesBatch(context.Context, []entities.Price) error
}

type portfolioRepository interface {
	Upsert(context.Context, entities.Portfolio) error
}

type tradesRepository interface {
	Delete(context.Context) error
	SaveTradesBatch(context.Context, []entities.Trade) error
}

type csvServices interface {
	GetPrices() []ports.PriceCSVDTO
	GetTrades() []ports.TradeCSVDTO
}

type ImportCSVDataInput struct{}

type ImportCSVDataOutput struct{}

type importCSVUC struct {
	csvReader           csvServices
	settingsRepository  settingsRepository
	assetsRepository    assetsRepository
	pricesRepository    pricesRepository
	portfolioRepository portfolioRepository
	tradesRepository    tradesRepository
	workers             int
}

func (s importCSVUC) Execute(ctx context.Context, _ ImportCSVDataInput) (ImportCSVDataOutput, error) {
	imported, err := s.settingsRepository.IsCsvImported(ctx)
	if err != nil {
		return ImportCSVDataOutput{}, err
	}

	if imported {
		return ImportCSVDataOutput{}, nil
	}

	assetA := entities.MustAsset(vos.MustID("84b3cf08-db0c-46aa-9a3c-cb97c98337ec"), "A")
	assetB := entities.MustAsset(vos.MustID("b94b8afb-dc8a-450f-8ebc-83cabe3b3c3a"), "B")
	if err := s.assetsRepository.Upsert(ctx, assetA); err != nil {
		return ImportCSVDataOutput{}, fmt.Errorf("%w:on upserting asset A", err)
	}

	if err := s.assetsRepository.Upsert(ctx, assetB); err != nil {
		return ImportCSVDataOutput{}, fmt.Errorf("%w:on upserting asset B", err)
	}

	s.pricesRepository.Delete(ctx)
	s.tradesRepository.Delete(ctx)

	if err := s.importPrices(ctx, importInput{
		assetA: assetA,
		assetB: assetB,
	}); err != nil {
		return ImportCSVDataOutput{}, err
	}

	if err := s.importTrades(ctx, importInput{
		assetA: assetA,
		assetB: assetB,
	}); err != nil {
		return ImportCSVDataOutput{}, err
	}

	alicePortID := vos.MustID("408186c6-b76a-4ad6-8d4a-9ace3762b997")
	portfolio, err := entities.NewPortfolio(
		alicePortID,
		vos.ParseToDecimal(100000),
		vos.ParseToDecimal(100000),
		map[string]entities.Position{},
	)
	if err != nil {
		return ImportCSVDataOutput{}, err
	}

	if err := s.portfolioRepository.Upsert(ctx, portfolio); err != nil {
		return ImportCSVDataOutput{}, fmt.Errorf("erro upsertando portfolio: %w", err)
	}

	if err := s.settingsRepository.SetCsvImported(ctx); err != nil {
		return ImportCSVDataOutput{}, fmt.Errorf("%w:on set csv to imported", err)
	}

	return ImportCSVDataOutput{}, nil
}

type importInput struct {
	assetA, assetB entities.Asset
}

func (s importCSVUC) importPrices(ctx context.Context, input importInput) error {
	pricesCSV := s.csvReader.GetPrices()
	prices := make([]entities.Price, 0, len(pricesCSV))

	for _, p := range pricesCSV {
		var assetID entities.AssetID
		switch p.Symbol {
		case "A":
			assetID = input.assetA.ID()
		case "B":
			assetID = input.assetB.ID()
		default:
			continue
		}

		priceEntity := entities.MustPrice(assetID, vos.ParseToDecimal(p.Price), p.Time.Time)
		prices = append(prices, priceEntity)
	}

	if err := s.pricesRepository.SavePricesBatch(ctx, prices); err != nil {
		return err
	}

	return nil
}

func (s importCSVUC) importTrades(ctx context.Context, input importInput) error {
	tradesCSV := s.csvReader.GetTrades()
	trades := make([]entities.Trade, 0, len(tradesCSV))

	for _, t := range tradesCSV {
		var asset entities.Asset
		switch t.Symbol {
		case "A":
			asset = input.assetA
		case "B":
			asset = input.assetB
		default:
			continue
		}

		side, err := vos.ParseSide(t.Side)
		if err != nil {
			return err
		}

		tradeEntity := entities.MustTrade(
			vos.NewID(),
			asset,
			side,
			vos.ParseToDecimal(t.Price),
			t.Quantity,
			t.Time.Time,
		)
		trades = append(trades, tradeEntity)
	}

	return s.tradesRepository.SaveTradesBatch(ctx, trades)
}

func NewImporterService(
	db *gbdb.Database,
	tx gbdb.Transactioner,
	csvReader csvServices,
	settingsRepository settingsRepository,
	assetsRepository assetsRepository,
	pricesRepository pricesRepository,
	portfolioRepository portfolioRepository,
	tradesRepository tradesRepository,
) gblib.UseCase[ImportCSVDataInput, ImportCSVDataOutput] {
	return gblib.New(
		importCSVUC{
			csvReader:           csvReader,
			settingsRepository:  settingsRepository,
			assetsRepository:    assetsRepository,
			pricesRepository:    pricesRepository,
			portfolioRepository: portfolioRepository,
			tradesRepository:    tradesRepository,
			workers:             5,
		},
		gblib.WithDB(db),
		gblib.WithTx(tx),
	)
}
