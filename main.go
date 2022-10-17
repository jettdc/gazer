package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Test test this thing on?")

	rand.Seed(time.Now().UnixNano())
	t, _ := time.ParseDuration(strconv.Itoa(rand.Intn(4)))
	time.Sleep(t)
	panic("ET")
}
