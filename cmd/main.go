package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gnur/slev"
)

func main() {
	sl, err := slev.Start(slev.UseDefaultHTTPServer("localhost:8821"))

	if err != nil {
		fmt.Println("Error:", err)
	}

	for i := 0; i < 10; i++ {

		id, err := sl.NewEvent("slevtest", "simpletest", map[string]int{"i": i})
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println(i, id)
	}

	time.Sleep(7 * time.Second)
	spew.Dump(sl.RawEvents())

	//sleep indefinitely
	ch := make(chan bool)
	<-ch
}
