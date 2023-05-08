package exchange

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
	"goweb/pkg/caculate"
	"goweb/pkg/console"
	"strconv"
	"time"
)

// 开合约 - 多
func BuySwapHandler(stop_loss float64, toke_profit float64) {
	console.Log(fmt.Sprintf("处理合约开多，读取传入参数，toke_profit(%.2f),stop_loss(%.2f)", toke_profit, stop_loss))

	// 设置杠杆倍数
	level, err := FuturesClient.NewChangeLeverageService().
		Symbol(CoinName).
		Leverage(leverage /*杠杆合约倍数*/).
		Do(context.Background())
	if err != nil {
		console.Log("初始化杠杆数失败")
		fmt.Println(err)
		return
	} else {
		console.Log(fmt.Sprintf("初始化 %s 杠杆数成功，当前杠杆为 %d", CoinName, level.Leverage))
	}

	for {
		res, err := FuturesClient.NewDepthService().
			Symbol(CoinName).
			Limit(5).
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		price, _ := strconv.ParseFloat(res.Asks[1].Price, 64)
		priceRound, _ := decimal.NewFromFloat(price).Round(swapPricePrecision /*金额精度*/).Float64()
		priceRound = priceRound + 1
		stopPrice := stop_loss     // 止损价格
		profitPrice := toke_profit // 止盈价格
		priceRoundStr := strconv.FormatFloat(priceRound, 'f', int(swapPricePrecision), 64)
		stopPriceRoundStr := strconv.FormatFloat(stopPrice, 'f', int(swapPricePrecision), 64)     // 止损价格
		profitPriceRoundStr := strconv.FormatFloat(profitPrice, 'f', int(swapPricePrecision), 64) // 止盈价格
		quantity := caculate.DivideFloat(float64(swapAmount), priceRound)
		quantity, _ = decimal.NewFromFloat(quantity).Round(swapPrecision /*币种购买数量精度*/).Float64()
		quantityStr := strconv.FormatFloat(quantity, 'f', int(swapPrecision), 64)

		order, err := FuturesClient.NewCreateOrderService().
			NewClientOrderID(swapRandoamId).
			Symbol(CoinName).
			PositionSide(futures.PositionSideTypeLong). //开多
			Side(futures.SideTypeBuy).
			Type(futures.OrderTypeLimit).
			TimeInForce(futures.TimeInForceTypeGTC).
			Quantity(quantityStr). // 这里需要自己计算，用 想购买的金额 / 限价价格
			Price(priceRoundStr).Do(context.Background())
		if err != nil {
			console.Log("下单失败")
			fmt.Println(err)
			return
		} else {
			console.Log(fmt.Sprintf("创建下单成功，订单id为: ", order.OrderID))

			//等 1 秒
			time.Sleep(1 * time.Second)

			// 获取所有订单，看是否存在未成交订单，存在则取消
			openOrders, err := FuturesClient.NewListOpenOrdersService().Symbol(CoinName).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if len(openOrders) > 0 {
				for _, curOpenOrder := range openOrders {
					// 存在未成交的订单，全部取消
					fmt.Println(curOpenOrder.OrderID)
					cancel, err := FuturesClient.NewCancelOrderService().Symbol(CoinName).
						OrderID(curOpenOrder.OrderID).Do(context.Background())
					if err != nil {
						console.Log("取消订单失败")
						fmt.Println(err)
						return
					} else {
						console.Log("取消订单成功")
						fmt.Println(cancel)
					}
				}
			} else {
				console.Log("合约开多 成功")

				//创建止损单
				order2, err := FuturesClient.NewCreateOrderService().
					Symbol(CoinName).
					Side(futures.SideTypeSell).
					ClosePosition(true).
					PositionSide(futures.PositionSideTypeLong). //开多
					Type(futures.OrderTypeStopMarket).
					TimeInForce(futures.TimeInForceTypeGTC).
					StopPrice(stopPriceRoundStr).
					Do(context.Background())
				if err != nil {
					console.Log("创建止损单失败")
					fmt.Println(err)
				} else {
					console.Log(fmt.Sprintf("创建止损单成功,订单id为：", order2.OrderID))
				}

				//创建止盈单
				order3, err := FuturesClient.NewCreateOrderService().
					Symbol(CoinName).
					Side(futures.SideTypeSell).
					ClosePosition(true).
					PositionSide(futures.PositionSideTypeLong). //开多
					Type(futures.OrderTypeTakeProfitMarket).
					TimeInForce(futures.TimeInForceTypeGTC).
					StopPrice(profitPriceRoundStr).
					Do(context.Background())
				if err != nil {
					console.Log("创建止盈单失败")
					fmt.Println(err)
				} else {
					console.Log("创建止盈单成功")
					console.Log(fmt.Sprintf("创建止盈单成功,订单id为：", order3.OrderID))
				}
				return
			}
		}
	}
}

// 开合约 - 空
func SellSwapHandler(stop_loss float64, toke_profit float64) {
	// 设置杠杆倍数
	level, err := FuturesClient.NewChangeLeverageService().
		Symbol(CoinName).
		Leverage(leverage /*杠杆合约倍数*/).
		Do(context.Background())
	if err != nil {
		console.Log("初始化杠杆数失败")
		fmt.Println(err)
		return
	} else {
		console.Log(fmt.Sprintf("初始化 %s 杠杆数成功，当前杠杆为 %d", CoinName, level.Leverage))
	}

	for {
		res, err := FuturesClient.NewDepthService().
			Symbol(CoinName).
			Limit(5).
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		price, _ := strconv.ParseFloat(res.Bids[1].Price, 64)
		priceRound, _ := decimal.NewFromFloat(price).Round(swapPricePrecision /*金额精度*/).Float64()
		priceRound = priceRound + 1
		stopPrice := stop_loss     // 止损价格
		profitPrice := toke_profit // 止盈价格
		priceRoundStr := strconv.FormatFloat(priceRound, 'f', int(swapPricePrecision), 64)
		stopPriceRoundStr := strconv.FormatFloat(stopPrice, 'f', int(swapPricePrecision), 64)     // 止损价格
		profitPriceRoundStr := strconv.FormatFloat(profitPrice, 'f', int(swapPricePrecision), 64) // 止盈价格
		quantity := caculate.DivideFloat(float64(swapAmount), priceRound)
		quantity, _ = decimal.NewFromFloat(quantity).Round(swapPrecision /*币种购买数量精度*/).Float64()
		quantityStr := strconv.FormatFloat(quantity, 'f', int(swapPrecision), 64)

		order, err := FuturesClient.NewCreateOrderService().
			NewClientOrderID(swapRandoamId).
			Symbol(CoinName).
			PositionSide(futures.PositionSideTypeShort). //开空
			Side(futures.SideTypeSell).
			Type(futures.OrderTypeLimit).
			TimeInForce(futures.TimeInForceTypeGTC).
			Quantity(quantityStr). // 这里需要自己计算，用 想购买的金额 / 限价价格
			Price(priceRoundStr).Do(context.Background())
		if err != nil {
			console.Log("下单失败")
			fmt.Println(err)
			return
		} else {
			console.Log(fmt.Sprintf("创建下单成功,订单id为：", order.OrderID))

			//等 1 秒
			time.Sleep(1 * time.Second)

			// 获取所有订单，看是否存在未成交订单，存在则取消
			openOrders, err := FuturesClient.NewListOpenOrdersService().Symbol(CoinName).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if len(openOrders) > 0 {
				for _, curOpenOrder := range openOrders {
					// 存在未成交的订单，全部取消
					fmt.Println(curOpenOrder.OrderID)
					cancel, err := FuturesClient.NewCancelOrderService().Symbol(CoinName).
						OrderID(curOpenOrder.OrderID).Do(context.Background())
					if err != nil {
						console.Log("取消订单失败") // 已经被交易
						fmt.Println(err)
					} else {
						console.Log("取消订单成功")
						fmt.Println(cancel)
					}
				}
			} else {
				console.Log("合约开空 成功")

				//创建止损单
				order2, err := FuturesClient.NewCreateOrderService().
					Symbol(CoinName).
					Side(futures.SideTypeBuy).
					ClosePosition(true).
					PositionSide(futures.PositionSideTypeShort). //开多
					Type(futures.OrderTypeStopMarket).
					TimeInForce(futures.TimeInForceTypeGTC).
					StopPrice(stopPriceRoundStr).
					Do(context.Background())
				if err != nil {
					console.Log("创建止损单失败")
					fmt.Println(err)
				} else {
					console.Log(fmt.Sprintf("创建止损单成功,订单id为：", order2.OrderID))
				}

				//创建止盈单
				order3, err := FuturesClient.NewCreateOrderService().
					Symbol(CoinName).
					Side(futures.SideTypeBuy).
					ClosePosition(true).
					PositionSide(futures.PositionSideTypeShort). //开多
					Type(futures.OrderTypeTakeProfitMarket).
					TimeInForce(futures.TimeInForceTypeGTC).
					StopPrice(profitPriceRoundStr).
					Do(context.Background())
				if err != nil {
					console.Log("创建止盈单失败")
					fmt.Println(err)
				} else {
					console.Log("创建止盈单成功")
					console.Log(fmt.Sprintf("创建止盈单成功,订单id为：", order3.OrderID))
				}
				return
			}
		}
	}
}
