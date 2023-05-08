## 拉取代码
```
先生成 Personal access tokens
然后用 git clone https://lbd2013:<Personal access tokens>@github.com/lbd2013/biance_deal 进行拉取
```

## 初始化
```
查看cmd命令：
go run main.go -h

创建表：
go run main.go migrate up

将k线数据拉取到本地：
go run main.go seed SeedKlinesTable

根据支撑线、压力线自动挂单：
go run . auto_trade
```

## 启动方法
```
开两个窗口
1. air (自动编译运行)
2. Mailhog (邮件服务器，用于注册)
```

## 链接
```
http://127.0.0.1:3307/api/v1/kline/toujianv7?coinName=BTCUSDT&startTime=2022-05-10&interval=1h&goType=down&timeWin=10
```