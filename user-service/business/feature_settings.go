package business

type FeatureSettings struct {
	MinHomeSquare float64
	MaxHomeSquare float64
}

func NewFeatureSettings() *FeatureSettings {
	return &FeatureSettings{
		MinHomeSquare: 20.0,
		MaxHomeSquare: 200.0,
	}
}
