package constants

// DelayType classifies a registered delay event.
type DelayType string

const (
	DelayWeather     DelayType = "WEATHER"
	DelayATC         DelayType = "ATC"
	DelayCatering    DelayType = "CATERING"
	DelayBaggage     DelayType = "BAGGAGE"
	DelayRefuel      DelayType = "REFUEL"
	DelayMaintenance DelayType = "MAINTENANCE"
	DelayOther       DelayType = "OTHER"
)

var DelayTypes = []DelayType{
	DelayWeather, DelayATC, DelayCatering, DelayBaggage, DelayRefuel, DelayMaintenance, DelayOther,
}

var DelayTypeText = map[DelayType]string{
	DelayWeather:     "天气",
	DelayATC:         "空管流控",
	DelayCatering:    "配餐延误",
	DelayBaggage:     "行李延误",
	DelayRefuel:      "加油延误",
	DelayMaintenance: "机务维护",
	DelayOther:       "其他",
}

func (t DelayType) Valid() bool {
	for _, candidate := range DelayTypes {
		if candidate == t {
			return true
		}
	}
	return false
}

func (t DelayType) Text() string { return DelayTypeText[t] }
