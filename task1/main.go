package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
)

func main() {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatalln("Fail to get time from ntp:", err)
	}

	fmt.Println(t)
}
