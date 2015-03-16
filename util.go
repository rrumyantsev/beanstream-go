package beanstream

import (
	//"fmt"
	"math/rand"
	"strconv"
	"time"
)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Util_randSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Util_randOrderId(num int) string {
	rnd := Util_randSeq(num)
	//fmt.Println("Timestamp: ", strconv.Itoa(int(time.Now().Unix())))
	rnd += strconv.Itoa(int(time.Now().Unix()))
	return rnd
}
