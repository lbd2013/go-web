package caculate

import (
	"github.com/shopspring/decimal"
	"time"
)

// AddFloat decimal类型加法
// return d1 + d2
func AddFloat(d1, d2 float64) float64 {
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Add(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}

// SubtractFloat decimal类型减法
// return d1 - d2
func SubtractFloat(d1, d2 float64) float64 {
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Sub(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}

// MultiplyFloat decimal类型乘法
// return d1 * d2
func MultiplyFloat(d1, d2 float64) float64 {
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Mul(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}

// DivideFloat decimal类型除法
// return d1 / d2
func DivideFloat(d1, d2 float64) float64 {
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Div(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}

func DatetimeToTimestamp(formatTimeStr string) (tm int64) {
	formatTime, _ := time.Parse("2006-01-02 15:04:05", formatTimeStr)
	return formatTime.Unix()
}

func TimeFromUnixNEscInt64(i int64) time.Time {
	return time.Unix(0, int64(i)*int64(time.Millisecond))
}

func StrTimeAdd(strTime string, addTime time.Duration) string {
	timeLayout := "2006-01-02 15:04:05"
	newTime, _ := time.Parse(timeLayout, strTime)
	newTime = newTime.Add(addTime)
	return newTime.UTC().Format(timeLayout)
}

//Calculate support pressure
// time1 对应的价格是是 price1，time2 对应的价格是 price2,targetTime 是你想计算的时间的价格
func SupportPressure(price1 float64, price2 float64, time1 string, time2 string, targetTime string) (targetPrice float64) {

	var timeLayout string = "2006-01-02 15:04:05"

	// 计算time时间跟第一个时间的间隔
	time1Parse, _ := time.Parse(timeLayout, time1)
	time1Unix := time1Parse.Unix()
	time2Parse, _ := time.Parse(timeLayout, time2)
	time2Unix := time2Parse.Unix()
	timeInterval := time2Unix - time1Unix

	// 计算目标时间到第一个时间的间隔
	targetParse, _ := time.Parse(timeLayout, targetTime)
	targetUnix := targetParse.Unix()
	targetInterval := targetUnix - time1Unix

	// 根据时间间隔比例，计算目标时间价格
	targetPrice = price1 - ((price1 - price2) * float64(targetInterval) / float64(timeInterval))

	return targetPrice
}
