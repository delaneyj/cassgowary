package cassgowary

import "math"

type Strength float64

var (
	Required = CreateStrengthWithDefaultWeight(1000, 1000, 1000)
	Strong   = CreateStrengthWithDefaultWeight(1, 0, 0)
	Medium   = CreateStrengthWithDefaultWeight(0, 1, 0)
	Weak     = CreateStrengthWithDefaultWeight(0, 0, 1)
)

func CreateStrength(a, b, c, w float64) Strength {
	result := 0.0
	result += math.Max(0, math.Min(1000, a*w)) * 1000000
	result += math.Max(0, math.Min(1000, b*w)) * 1000
	result += math.Max(0, math.Min(1000, c*w))
	return Strength(result)
}

func CreateStrengthWithDefaultWeight(a, b, c float64) Strength {
	return CreateStrength(a, b, c, 1)
}

func ClipStrength(value Strength) Strength {
	r, v := float64(Required), float64(value)
	return Strength(math.Max(0, math.Min(r, v)))
}

func (s Strength) Float() Float {
	return Float(s)
}
