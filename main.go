package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"inter/store2"
)

func main() {
	now := time.Now()

	defer func() {
		fmt.Printf("duration: %v\n", time.Since(now))
	}()

	store := store2.New() // change here the implementation, and here we go
	N := 5000000
	myarr := make([]string, N)

	for i := 0; i < N; i++ {
		myarr[i] = uuid.New().String()
	}

	go func() {
		for _, key := range myarr {
			store.Add(key, "ok")
		}

		fmt.Println("Add done")
	}()

	time.Sleep(time.Second)
	fmt.Println("start")

	for i := 0; i < len(myarr); i++ {
		key := myarr[i]
		value, err := store.Get(key)
		if err != nil {
			fmt.Printf("i = %d, value = %s :%v\n", i, key, err)

			i = 0
			time.Sleep(time.Millisecond * 100)
			continue
		}

		if value != "ok" {
			panic("not OK")
		}
	}

	fmt.Println("Done")
}
