package entities

type Position struct {
	assetID  AssetID
	quantity int
}

func (p Position) AssetID() AssetID {
	return p.assetID
}

func (p Position) Quantity() int {
	return p.quantity
}

func NewPosition(
	assetID AssetID,
	quantity int,
) (Position, error) {
	if assetID.IsZero() {
		return Position{}, ErrInvalidAssetID
	}

	return Position{
		assetID:  assetID,
		quantity: quantity,
	}, nil
}
