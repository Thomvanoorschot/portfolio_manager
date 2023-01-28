package enums

type AssetType string

const (
	AssetTypeUnknown AssetType = "Unknown"
	Equity                     = "Equity"
	ETF                        = "ETF"
	Cash                       = "Cash"
)
