package main

import (
	"fmt"
	"go-cache"
	"time"
)

type User struct {
	Name string
	Age  int
}

func main() {
	c := cache.NewCache[string, User]()
	u := &User{Name: "Mike", Age: 20}
	c.Set("key1", u, 0)

	if val1, ok := c.Get("key1"); ok {
		fmt.Printf("key1: %v\n", val1)
	} else {
		fmt.Println("none")
	}

	time.Sleep(0 * time.Second)

	if val1, ok := c.Get("key1"); ok {
		fmt.Printf("key1: %v\n", val1)
	} else {
		fmt.Println("none")
	}
}
