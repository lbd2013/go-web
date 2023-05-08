package cmd

import (
	"context"
	"fmt"
	godump "github.com/favframework/debug"
	"goweb/pkg/config"
	"goweb/pkg/console"
	"goweb/pkg/exchange"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var CmdAutoTrade = &cobra.Command{
	Use:   "auto_trade",
	Short: "auto trade between pressure price and support price",
	Run:   runAutoTrade,
	Args:  cobra.NoArgs, // 不允许传参
}

// 计算小时线趋势
func caculateHourTo(checkTime string) (float64, float64) {
	var timeLayout string = "2006-01-02 15:04:05"
	var intervalSec int64 = 60 * 60 // 用来计算两个时间之间间隔多少个小时
	var pressurePrice float64
	var supportPrice float64

	// 读取配置
	// 两个最高价格的时间点以及价格
	highest1Time := config.GetString("HIGHEST_TIME_1")
	highest1Price := config.GetFloat64("HIGHEST_PRICE_1")
	highest2Time := config.GetString("HIGHEST_TIME_2")
	highest2Price := config.GetFloat64("HIGHEST_PRICE_2")

	// 两个最低价格的时间点以及价格
	lowest1Time := config.GetString("LOWEST_TIME_1")
	lowest1Price := config.GetFloat64("LOWEST_PRICE_1")
	lowest2Time := config.GetString("LOWEST_TIME_2")
	lowest2Price := config.GetFloat64("LOWEST_PRICE_2")

	console.Log(fmt.Sprintf("读取配置价格成功"))
	console.Log(fmt.Sprintf("两个最高价分别为：%.6f(%s)\t %.6f(%s)\t ", highest1Price, highest1Time, highest2Price, highest2Time))
	console.Log(fmt.Sprintf("两个最低价分别为：%.6f(%s)\t %.6f(%s)\t", lowest1Price, lowest1Time, lowest2Price, lowest2Time))

	// 计算两个最高点时间的间隔
	highest1TimeParse, _ := time.Parse(timeLayout, highest1Time)
	highest1TimeUnix := highest1TimeParse.Unix()
	highest2TimeParse, _ := time.Parse(timeLayout, highest2Time)
	highest2TimeUnix := highest2TimeParse.Unix()
	highestInterval := (highest2TimeUnix - highest1TimeUnix) / intervalSec
	console.Log(fmt.Sprintf("两个最高点时间间隔 %d 个小时", highestInterval))

	// 计算要计算时间跟前一个时间的间隔
	checkTimeParse, _ := time.Parse(timeLayout, checkTime)
	checkTimeUnix := checkTimeParse.Unix()
	checkTimeHighInterval := (checkTimeUnix - highest1TimeUnix) / intervalSec
	console.Log(fmt.Sprintf("要计算时间跟第 1 个最高点时间间隔 %d 个小时", checkTimeHighInterval))

	// 计算两个最低点时间的间隔
	lowest1TimeParse, _ := time.Parse(timeLayout, lowest1Time)
	lowest1TimeUnix := lowest1TimeParse.Unix()
	lowest2TimeParse, _ := time.Parse(timeLayout, lowest2Time)
	lowest2TimeUnix := lowest2TimeParse.Unix()
	lowestInterval := (lowest2TimeUnix - lowest1TimeUnix) / intervalSec
	console.Log(fmt.Sprintf("两个最低点时间的间隔 %d 个小时", lowestInterval))

	// 计算要计算时间跟前一个时间的间隔
	checkTimeLowInterval := (checkTimeUnix - lowest1TimeUnix) / intervalSec
	console.Log(fmt.Sprintf("要计算时间跟第 1 个最低点时间间隔 %d 个小时", checkTimeLowInterval))

	// 根据比例计算出压力位
	pressurePrice = highest1Price + (highest2Price-highest1Price)/float64(highestInterval)*float64(checkTimeHighInterval)

	// 根据比例计算出支撑位
	supportPrice = lowest1Price + (lowest2Price-lowest1Price)/float64(lowestInterval)*float64(checkTimeLowInterval)

	console.Log(fmt.Sprintf("两个最高价分别为：%.6f(%s)\t %.6f(%s)\t 差值为:%.6f\t 两个最高价间隔 %d 个小时\t 要计算的时间跟第 1 个最高价间隔 %d 个小时\t 计算出来的当前时间的压力价为:%.6f\t",
		highest1Price, highest1Time, highest2Price, highest2Time, highest2Price-highest1Price, highestInterval, checkTimeHighInterval, pressurePrice))
	console.Log(fmt.Sprintf("两个最低价分别为：%.6f(%s)\t %.6f(%s)\t 差值为:%.6f\t 两个最低价间隔 %d 个小时\t 要计算的时间跟第 1 个最低价的间隔 %d 个小时\t 计算出来的当前时间的支撑价为:%.6f\t",
		lowest1Price, lowest1Time, lowest2Price, lowest2Time, lowest2Price-lowest1Price, lowestInterval, checkTimeLowInterval, supportPrice))
	return pressurePrice, supportPrice
}

func runAutoTrade(cmd *cobra.Command, args []string) {
	console.Log("自动交易开始...")
	var timeLayout string = "2006-01-02 15:04:05"

	// 1. 查看账户持仓，如果有持仓，退出
	res, err := exchange.FuturesClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		panic(err)
	}

	console.Log("以USD计价的所需起始保证金总额为：" + res.TotalInitialMargin)
	totalInitialMargin, _ := strconv.ParseFloat(res.TotalInitialMargin, 64)
	if totalInitialMargin > 0 {
		console.Log("始保证金总额大于0， 存在未平合约，本次结束")
		return
	} else {
		console.Log("始保证金总额为0，没有未平合约，继续")
	}

	// 2. 获取所有挂单并取消
	console.Log("获取所有挂单并取消...")
	openOrders, err := exchange.FuturesClient.NewListOpenOrdersService().Symbol(exchange.CoinName).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	for _, o := range openOrders {
		godump.Dump("订单id:" + strconv.Itoa(int(o.OrderID)))
		_, err := exchange.FuturesClient.NewCancelOrderService().Symbol(exchange.CoinName).
			OrderID(o.OrderID).Do(context.Background())
		if err != nil {
			panic(err)
		} else {
			console.Log("取消成功")
		}
	}

	// 3. 获取当前价格、当前价格时间
	klines, err := exchange.FuturesClient.NewKlinesService().Symbol(exchange.CoinName).
		Interval("1h").
		Limit(1).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	currPrice, _ := strconv.ParseFloat(klines[0].Close, 64)

	// 3. 根据时间获取当前压力线，支撑线
	end_time := time.Unix(klines[0].OpenTime/1000, 0)
	end_time_Format := end_time.Format(timeLayout)
	//end_time_Format = "2022-08-25 08:00:00"
	pressurePrice, supportPrice := caculateHourTo(end_time_Format)
	console.Success(fmt.Sprintf("当前价为%.6f\t 压力价为%.6f\t 支撑价位%.6f\t", currPrice, pressurePrice, supportPrice))

	canDo, flag, stop_loss_price, take_profit_price := validBuyOrSell(currPrice, supportPrice, pressurePrice)
	console.Log(fmt.Sprintf("根据当前价格、压力位、支撑位、最大允许止损金额判断是开空还是开多结果：canDo(%t), flag(%s), stop_loss_price(%.6f), take_profit_price(%.6f)", canDo, flag, stop_loss_price, take_profit_price))
	if canDo == true {
		if flag == "buy" {
			exchange.BuySwapHandler(stop_loss_price, take_profit_price)
		} else if flag == "sell" {
			exchange.SellSwapHandler(stop_loss_price, take_profit_price)
		}
	}
	return
}

// 根据当前价格、压力位、支撑位、最大允许止损金额判断是开空还是开多
func validBuyOrSell(currPrice float64, supportPrice float64, pressurePrice float64) (bool, string, float64, float64) {
	buy_toke_profit := config.GetFloat64("BUY_TAKE_PROFIT_PRICE")
	sell_toke_profit := config.GetFloat64("SELL_TAKE_PROFIT_PRICE")
	stop_loss := config.GetFloat64("STOP_LOSS_PRICE")
	max_stop_loss := config.GetFloat64("MAX_STOP_LOSS")
	swap_amount := config.GetFloat64("SWAP_AMOUNT")
	stop_loss_rate := config.GetFloat64("STOP_LOSS_RATE")

	console.Log(fmt.Sprintf("读取配置参数，buy_toke_profit(%.6f),sell_toke_profit(%.6f),stop_loss(%.6f),max_stop_loss(%.6f),swap_amount(%.6f),stop_loss_rate(%.6f)", buy_toke_profit, sell_toke_profit, stop_loss, max_stop_loss, swap_amount, stop_loss_rate))

	// 判断当前价格是否处于支撑价跟压力价之间，如果时，不下单
	if currPrice >= supportPrice && currPrice <= pressurePrice {
		console.Log("当前价格位于支撑价跟压力价之间，不下单")
		return false, "", 0, 0
	}

	if buy_toke_profit == -1 || sell_toke_profit == -1 {
		console.Exit(".env 文件中 BUY_TAKE_PROFIT_PRICE / SELL_TAKE_PROFIT_PRICE 值为 -1，请修改")
	}

	// 如果当前价格大于压力价，正常情况下是开多
	if currPrice > pressurePrice {
		// 1. 获取止损价格
		if stop_loss == -1 { // 止损价设置为 -1 ，多单使用被突破的压力位作为止损价格， 否则，使用配置中的止损价作为止损价格
			stop_loss = pressurePrice
		}

		// 2. 根据开单价（当前价格）、止损价格，计算止损金额会不会超过最大允许止损金额，如果小于，可以开多
		willStopLoss := (currPrice - stop_loss) / currPrice * swap_amount
		if willStopLoss <= max_stop_loss {
			console.Log(fmt.Sprintf("止损价为(%.6f), 当前价为(%.6f)，已突破压力位, 止损时，将会损失金额(%.6f)未超出最大允许止损金额(%.6f)，可以开多", stop_loss, currPrice, willStopLoss, max_stop_loss))
			return true, "buy", stop_loss, buy_toke_profit
		} else {
			console.Log(fmt.Sprintf("止损价为(%.6f), 当前价为(%.6f)，已突破压力位, 止损时，将会损失金额(%.6f),超出最大允许止损金额(%.6f)，不可以开多。", stop_loss, currPrice, willStopLoss, max_stop_loss))
			return false, "", 0, 0
		}

		//console.Log("判断能否开空...")
		////3. 第二步中算出来的止损金额超出了最大止损金额，说明现在未到三角形的尖尖部分，还处于三角形腹部，判断能否开空
		////3.1 止损比例为 .env文件中的 STOP_LOSS_RATE 值
		////3.2 根据止损比例计算止损时损失金额，对比该金额跟止盈金额，如果小于止盈金额，不下单
		////3.3 使用支撑线作为止盈点
		//take_profit_price := supportPrice
		//willStopLoss = swap_amount * stop_loss_rate
		//willTakeProfit := (currPrice - take_profit_price) / currPrice * swap_amount
		//
		//if willTakeProfit <= willStopLoss {
		//	console.Log(fmt.Sprintf("打算开单价格为(%.6f),止损比例为(%.6f),止损时损失金额为(%.6f),止盈时利润为(%.6f),止盈利润小于等于止损损失，不可开空", currPrice, stop_loss_rate, willStopLoss, willTakeProfit))
		//	return false, "", 0, 0
		//} else {
		//	// 判断当前价格是否超出压力价 * (1+stop_loss_rate)，如果是，不开空
		//	pressureRatePrice := pressurePrice * (1 + stop_loss_rate)
		//	stop_loss_price := currPrice * (1 + stop_loss_rate)
		//	console.Log(fmt.Sprintf("当前价格为(%.6f),压力位为(%.6f),压力价 * (1+止损比例)为(%.6f)", currPrice, pressurePrice, pressureRatePrice))
		//	if currPrice > pressureRatePrice {
		//		console.Log("当前价已经超出压力位太多，不可以开空")
		//		return false, "", 0, 0
		//	} else {
		//		console.Log("当前价没超出压力位太多，可以开空")
		//	}
		//
		//	console.Log(fmt.Sprintf("打算开单价格为(%.6f), 止损价为(%.6f), 止盈价为(%.6f),止损比例为(%.6f),止损时损失金额为(%.6f),止盈时利润为(%.6f),止盈利润大于止损损失，可以开空", currPrice, stop_loss_price, take_profit_price, stop_loss_rate, willStopLoss, willTakeProfit))
		//	return true, "sell", stop_loss_price, take_profit_price
		//}
	} else if currPrice < supportPrice { // 如果当前价格小于支撑，正常是开空
		// 1. 获取止损价格
		if stop_loss == -1 { // 止损价设置为 -1 ，空单使用被跌破的支撑位作为止损价格， 否则，使用配置中的止损价作为止损价格
			stop_loss = supportPrice
		}

		// 2. 根据开单价（当前价格）、止损价格，计算止损金额会不会超过最大允许止损金额，如果小于，可以开空
		willStopLoss := (stop_loss - currPrice) / currPrice * swap_amount
		if willStopLoss <= max_stop_loss {
			console.Log(fmt.Sprintf("止损价为(%.6f), 当前价为(%.6f)，已跌破, 止损时，将会损失金额(%.6f)未超出最大允许止损金额(%.6f)，可以开空", stop_loss, currPrice, willStopLoss, max_stop_loss))
			return true, "sell", stop_loss, sell_toke_profit
		} else {
			console.Log(fmt.Sprintf("止损价为(%.6f), 当前价为(%.6f)，已跌破, 止损时，将会损失金额(%.6f),超出最大允许止损金额(%.6f)，不可以开空。", stop_loss, currPrice, willStopLoss, max_stop_loss))
			return false, "", 0, 0
		}

		//console.Log("判断能否开多...")
		////3. 第二步中算出来的止损金额超出了最大止损金额，说明现在未到三角形的尖尖部分，还处于三角形腹部，判断能否开多
		////3.1 止损比例为 .env文件中的 STOP_LOSS_RATE 值
		////3.2 根据止损比例计算止损时损失金额，对比该金额跟止盈金额，如果小于止盈金额，不下单
		////3.3 使用压力位作为止盈点
		//take_profit_price := pressurePrice
		//willStopLoss = swap_amount * stop_loss_rate
		//willTakeProfit := (take_profit_price - currPrice) / currPrice * swap_amount
		//
		//if willTakeProfit <= willStopLoss {
		//	console.Log(fmt.Sprintf("打算开单价格为(%.6f),止损比例为(%.6f),止损时损失金额为(%.6f),止盈时利润为(%.6f),止盈利润小于等于止损损失，不可开多", currPrice, stop_loss_rate, willStopLoss, willTakeProfit))
		//	return false, "", 0, 0
		//} else {
		//	// 判断当前价格是否跌超出支撑价 * (1- stop_loss_rate)，如果是，不开多
		//	supportRatePrice := supportPrice * (1 - stop_loss_rate)
		//	stop_loss_price := currPrice * (1 - stop_loss_rate)
		//	console.Log(fmt.Sprintf("当前价格为(%.6f),支撑价为(%.6f),支撑价 * (1-止损比例)为(%.6f)", currPrice, supportPrice, supportRatePrice))
		//	if currPrice < supportRatePrice {
		//		console.Log("当前价已经跌超出支撑位太多，不可以开多")
		//		return false, "", 0, 0
		//	} else {
		//		console.Log("当前价没跌超出支撑位太多，可以开多")
		//	}
		//
		//	console.Log(fmt.Sprintf("打算开单价格为(%.6f), 止损价为(%.6f), 止盈价为(%.6f),止损比例为(%.6f),止损时损失金额为(%.6f),止盈时利润为(%.6f),止盈利润大于止损损失，可以开空", currPrice, stop_loss_price, take_profit_price, stop_loss_rate, willStopLoss, willTakeProfit))
		//	return true, "buy", stop_loss_price, take_profit_price
		//}
	} else {
		console.Log("啥也不是，不能开单")
		return false, "", 0, 0
	}
}
