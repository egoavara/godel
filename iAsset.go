package godel

type Asset interface {
	Name() string
	Type() AssetType
}
type AssetType uint8

const (
	AssetSkeleton  AssetType = iota
	AssetSkin      AssetType = iota
	AssetAnimation AssetType = iota
)
