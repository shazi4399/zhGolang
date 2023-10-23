package main

import "fmt"

type Phone interface {
	call()
	seenMessage()
}

type Huawei struct {
	name  string
	price int64
}

type Xiaomi struct {
	name  string
	price float64
}

func (huawei Huawei) call() {
	fmt.Printf("%s 有打电话功能.....\n", huawei.name)
}

func (huawei Huawei) seenMessage() {
	fmt.Printf("%s 有发短信功能.....\n", huawei.name)
}
func main() {
	mate30 := Huawei{
		name:  "Mate 30",
		price: 6999,
	}
	mate30.call()
	mate30.seenMessage()
	//通过new来检测是否实现了接口，如果没有实现会报错。
	var _ Phone = new(Huawei)
	fmt.Println()
}
