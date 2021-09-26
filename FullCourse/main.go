package main

import (
	"fmt"
	"sync"
)
var wg=sync.WaitGroup{}
func main()  {
	/* A Goroutine (function prefixed with go) is a function or method which executes independently and simultaneously in connection with any other Goroutines present in your program. 
	Or in other words, every concurrently executing activity in Go language is known as a Goroutines. You can consider a Goroutine like a light weighted thread. */
	msg :="Hello"
	wg.Add(1)
	// Annonymous function
	go func ()  {
		fmt.Println(msg)
		wg.Done()
	}()
	msg="Hi Go!!!"
	fmt.Println(msg)
	wg.Wait()
	// time.Sleep(1000*time.Microsecond)		// Hello\n Hi Go!!!\n 
}