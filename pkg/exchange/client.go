package exchange

import (
	"context"
	"crypto/tls"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"goweb/pkg/config"
	"net/http"
)

// 币安秘钥
var apiKey string
var secretKey string

// 要操作的币种
var CoinName string

// 币安操作对象，Client是现货，FuturesClient是合约
var Client *binance.Client
var FuturesClient *futures.Client

// 现货配置
var spotAmount int           // 每次下单金额,最小10
var spotPrecision int32      // 现货下单数量精度
var spotPricePrecision int32 // 现货金额精度

// 合约配置
var swapAmount int           // 每次下单金额,最小 10*杠杆倍数
var swapPrecision int32      // 合约下单数量精度
var swapPricePrecision int32 // 合约金额精度
var leverage int             // 合约杠杆倍数

// 现货随机id，避免重复下单,只能避免多线程重复下单，如果是单线程的话，此配置没用
var spotRandoamId string = "sdajfxcjvhxzjhvjkhsadf"

// 合约随机id，避免重复下单,只能避免多线程重复下单，如果是单线程的话，此配置没用
var swapRandoamId string = "xcmvnabsdmbxncvbisadif"

func Init() {
	CoinName = config.GetString("CoinName")                     // 币种名称
	apiKey = config.GetString("BINANCE_API_KEY")                // 币安秘钥
	secretKey = config.GetString("BINANCE_SECRET_KEY")          // 币安秘钥
	spotAmount = config.GetInt("SPOT_AMOUNT")                   // 每次下单金额,最小10
	spotPrecision = config.GetInt32("SPOT_PRECISION")           // 现货下单数量精度
	spotPricePrecision = config.GetInt32("SPOT_PRICEPRECISION") // 现货金额精度
	swapAmount = config.GetInt("SWAP_AMOUNT")                   // 每次下单金额,最小 10*杠杆倍数
	swapPrecision = config.GetInt32("SWAP_PRECISION")           // 合约下单数量精度
	swapPricePrecision = config.GetInt32("SWAP_PRICEPRECISION") // 合约金额精度
	leverage = config.GetInt("SWAP_LEVERAGE")                   // 合约杠杆倍数

	// 设置代理
	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	http.DefaultClient.Transport = tr

	// 初始化现货跟合约操作对象
	Client = binance.NewClient(apiKey, secretKey)
	Client.NewSetServerTimeService().Do(context.Background())
	FuturesClient = binance.NewFuturesClient(apiKey, secretKey) // USDT-M Futures
	FuturesClient.NewSetServerTimeService().Do(context.Background())
}
