package engine

import (
	"fmt"
	"matching/enum"
	"time"

	"github.com/shopspring/decimal"
)

func (o *OrderBookPriority) GetAskDepth(size int) [][2]string {
	return o.depth(o.sellBook, size)
}

func (o *OrderBookPriority) GetBidDepth(size int) [][2]string {
	return o.depth(o.buyBook, size)
}

func (o *OrderBookPriority) depth(queue *OrderQueue, size int) [][2]string {
	queue.Lock()
	defer queue.Unlock()

	max := len(queue.depth)
	if size <= 0 || size > max {
		size = max
	}

	return queue.depth[0:size]
}

func (o *OrderBookPriority) depthTicker(que *OrderQueue) {

	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)

	for {
		<-ticker.C
		func() {
			o.w.Lock()
			defer o.w.Unlock()

			que.Lock()
			defer que.Unlock()
			que.depth = [][2]string{}
			depthMap := make(map[string]string)

			if que.pq.Len() > 0 {

				for i := 0; i < que.pq.Len(); i++ {
					item := (*que.pq)[i]

					price := FormatDecimal2String(item.GetPrice(), o.priceDigit)

					if _, ok := depthMap[price]; !ok {
						depthMap[price] = FormatDecimal2String(item.GetQuantity(), o.quantityDigit)
					} else {
						old_qunantity, _ := decimal.NewFromString(depthMap[price])
						depthMap[price] = FormatDecimal2String(old_qunantity.Add(item.GetQuantity()), o.quantityDigit)
					}
				}

				//按价格排序map
				que.depth = sortMap2Slice(depthMap, que.Top().GetOrderSide())
			}
		}()
	}
}

func FormatDecimal2String(d decimal.Decimal, digit int) string {
	f, _ := d.Float64()
	format := "%." + fmt.Sprintf("%d", digit) + "f"
	return fmt.Sprintf(format, f)
}

func quickSort(nums []string, asc_desc string) []string {
	if len(nums) <= 1 {
		return nums
	}

	spilt := nums[0]
	left := []string{}
	right := []string{}
	mid := []string{}

	for _, v := range nums {
		vv, _ := decimal.NewFromString(v)
		sp, _ := decimal.NewFromString(spilt)
		if vv.Cmp(sp) == -1 {
			left = append(left, v)
		} else if vv.Cmp(sp) == 1 {
			right = append(right, v)
		} else {
			mid = append(mid, v)
		}
	}

	left = quickSort(left, asc_desc)
	right = quickSort(right, asc_desc)

	if asc_desc == "asc" {
		return append(append(left, mid...), right...)
	} else {
		return append(append(right, mid...), left...)
	}

	//return append(append(left, mid...), right...)
}

func sortMap2Slice(m map[string]string, ask_bid enum.OrderSide) [][2]string {
	res := [][2]string{}
	keys := []string{}
	for k, _ := range m {
		keys = append(keys, k)
	}

	if ask_bid == enum.SideSell {
		keys = quickSort(keys, "asc")
	} else {
		keys = quickSort(keys, "desc")
	}

	for _, k := range keys {
		res = append(res, [2]string{k, m[k]})
	}
	return res
}
