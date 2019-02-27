package main

import (
	"github.com/adshao/go-binance"
	"context"
	"fmt"
	"strconv"
	"time"
	"math"
)



var (
	apiKey = ""
	secretKey = ""

)

type Accuracy func() float64

func (this Accuracy) Equal(a, b float64) bool {
	return math.Abs(a-b) < this()
}

func (this Accuracy) Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > this()
}

func (this Accuracy) Smaller(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > this()
}

func (this Accuracy) GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < this()
}

func (this Accuracy) SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < this()
}



var client = binance.NewClient(apiKey, secretKey)

func sysmbol_price(sysmbol string) (price string){
	prices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, pair := range prices{
		if pair.Symbol == sysmbol{
			return pair.Price
		}
	}
	return ""
}

func get_balance(coin string) (qty float64) {
	res, err := client.NewGetAccountService().Do(context.Background())
	var count float64
	if err != nil {
		fmt.Println(err)
		return
	}
	for _,each := range res.Balances{
		if each.Asset == coin{
			count,_ = strconv.ParseFloat(each.Free,64)
		}
	}
	return count

}

func trend() (state bool)  {
	endTime := time.Now().UnixNano()/1e6
	startTime := time.Now().Add(-time.Minute*60).UnixNano()/1e6
	fmt.Println(endTime,startTime)
	trades, err := client.NewAggTradesService().
		Symbol("ETHBTC").StartTime(startTime).
		EndTime(endTime).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	var start_temp float64 = 0
	for _, t := range trades[:10] {
		t_Price,_ := strconv.ParseFloat(t.Price,64)
		start_temp = t_Price + start_temp
	}
	start_ave := start_temp/10

	var end_temp float64 = 0
	for _, t := range trades[len(trades)-10:] {
		t_Price,_ := strconv.ParseFloat(t.Price,64)
		end_temp = t_Price + end_temp
	}
	end_ave := end_temp/10
	fmt.Println(start_ave,end_ave)

	var a Accuracy = func() float64 { return 0.0000001 }

	if a.Greater(start_temp,end_ave){
		return false
	}else {
		return true
	}

}

func trend_small() (state int){
	endTime := time.Now().UnixNano()/1e6
	startTime := time.Now().Add(-time.Minute*10).UnixNano()/1e6
	fmt.Println(endTime,startTime)
	trades, err := client.NewAggTradesService().
		Symbol("ETHBTC").StartTime(startTime).
		EndTime(endTime).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	var start_temp float64 = 0
	for _, t := range trades[:10] {
		t_Price,_ := strconv.ParseFloat(t.Price,64)
		start_temp = t_Price + start_temp
	}
	start_ave := start_temp/10

	var end_temp float64 = 0
	for _, t := range trades[len(trades)-10:] {
		t_Price,_ := strconv.ParseFloat(t.Price,64)
		end_temp = t_Price + end_temp
	}
	end_ave := end_temp/10
	fmt.Println(start_ave,end_ave)

	var a Accuracy = func() float64 { return 0.0000001 }

	if a.Greater(start_ave,end_ave){
		fmt.Println(-1)
		return -1
	}else if (a.Equal(start_ave,end_ave)) {
		fmt.Println(0)
		return 0
	}else if(a.Smaller(start_ave,end_ave)){
		fmt.Println(1)
		return 1
	} else {
		return
	}
}


func Round2(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}

func Max_min_arry(arr []float64) (float64,float64){
	max := arr[0]
	min := arr[0]
	var a Accuracy = func() float64 { return 0.0000001 }

	for _, t := range arr{
		if a.Greater(t,max) {
			max = t
		}
		if a.Smaller(t,min){
			min = t
		}
	}
	return max, min
	}


func limit_order(sell_per float64, buy_per float64)  {
	res, err := client.NewDepthService().Symbol("ETHBTC").
		Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}
	//market_price := sysmbol_price("ETHBTC")
	buy_price := res.Bids[0].Price
	sell_price := res.Bids[0].Price

	f_sell_price,_ := strconv.ParseFloat(sell_price,64)
	fmt.Println(f_sell_price)
	f_sell_price = Round2(f_sell_price*sell_per,6)
	f_buy_price,_ := strconv.ParseFloat(buy_price,64)
	fmt.Println(f_buy_price)
	f_buy_price = Round2(f_buy_price*buy_per,6)

	s_sell_price := strconv.FormatFloat(f_sell_price,'g',-1,64)
	s_buy_price := strconv.FormatFloat(f_buy_price,'g',-1,64)
	fmt.Println(s_sell_price,s_buy_price)

	TradeSize  := get_balance("ETH") *0.98
	TradeSize_count := Round2(TradeSize,3)
	s_tradesize := strconv.FormatFloat(TradeSize_count,'g',-1,64)
	fmt.Println(s_tradesize)





	//sell_trade_size := strconv.FormatFloat(Tradesize,'g',-1,64)
	if ((f_sell_price-f_buy_price)>0) && (f_buy_price < 0.00031){

		//endTime := time.Now().UnixNano()/1e6
		//startTime := time.Now().Add(-time.Hour*1).UnixNano()/1e6
		//trades, err := client.NewAggTradesService().
		//	Symbol("ETHBTC").StartTime(startTime).
		//	EndTime(endTime).Do(context.Background())
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//var prices []float64
		//for _, t := range trades{
		//	t_Price,_ := strconv.ParseFloat(t.Price,64)
		//	prices = append(prices,t_Price)
		//}

		//max_price , min_price := Max_min_arry(prices)
		//var a Accuracy = func() float64 { return 0.0000001 }
		//if (a.Greater(max_price*1.008,f_sell_price)) && (a.Smaller(min_price*0.992,f_buy_price)){
			_, err1 := client.NewCreateOrderService().Symbol("ETHBTC").Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).TimeInForce(binance.TimeInForceGTC).Quantity(s_tradesize).Price(s_sell_price).Do(context.Background())
			if err1 != nil {
				fmt.Println(err1)
				return
			}

			_, err2 := client.NewCreateOrderService().Symbol("ETHBTC").Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).TimeInForce(binance.TimeInForceGTC).Quantity(s_tradesize).Price(s_buy_price).Do(context.Background())
			if err2 != nil {
				fmt.Println(err2)
				return
			}
		//}
	}
}

func order_state()  (state bool){
	orders, err := client.NewListOpenOrdersService().Symbol("ETHBTC").
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, o := range orders {
		fmt.Println(o)
		return true
	}
	return false
}





func main()  {
	var counts int = 0
	for {
		if order_state() == false{
			counts =0
			time.Sleep(time.Minute*60*5)

			if trend() == false {
				limit_order(1.001, 0.995)
			}
			//} else {
			//	limit_order(1.010,1)
			//}
		}
		counts ++
		time.Sleep(time.Minute*1)
		if (counts >5440){
			orders, _ := client.NewListOpenOrdersService().Symbol("ETHBTC").
				Do(context.Background())
			for _, o := range orders{
				Order_id := o.OrderID
				Order_Side := o.Side
				Tradesize := o.OrigQuantity
				if (Order_Side == "SELL") {
					if (trend_small() == -1) {
						_, err := client.NewCancelOrderService().Symbol("ETHBTC").
							OrderID(Order_id).Do(context.Background())
						if err == nil {
							client.NewCreateOrderService().Symbol("ETHBTC").Side(binance.SideTypeSell).
								Type(binance.OrderTypeMarket).Quantity(Tradesize).Do(context.Background())
						}
					}
				} else if (Order_Side == "BUY") {

					if (trend_small() == 1) {
						_, err := client.NewCancelOrderService().Symbol("ETHBTC").
							OrderID(Order_id).Do(context.Background())
						if err == nil {
							client.NewCreateOrderService().Symbol("ETHBTC").Side(binance.SideTypeBuy).
								Type(binance.OrderTypeMarket).Quantity(Tradesize).Do(context.Background())
						}
					}
				}
			}
		}
	}
}
