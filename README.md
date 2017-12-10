# go-epa

[![Build Status](https://travis-ci.org/minchao/go-epa.svg?branch=master)](https://travis-ci.org/minchao/go-epa)
[![codecov](https://codecov.io/gh/minchao/go-epa/branch/master/graph/badge.svg)](https://codecov.io/gh/minchao/go-epa)

台灣環保署開放資料 Golang library

## 安裝

使用 `go get` 指令安裝

```
go get github.com/minchao/go-epa
```

## 使用

```go
import "github.com/minchao/go-epa"
```

首先初始化 API Client，接著就可以用它來存取環保署的開放資料 API，如範例：

```go
client := epa.NewClient("TOKEN", nil)
ctx := context.Background()

// 取得空氣品質預測資料
options := url.Values{}
options.Set("sort", "PublishTime")
credit, err := client.GetAirQualityForecast(ctx, options)
```

## License

This library is distributed under the BSD-style license found in the [LICENSE](./LICENSE) file.
