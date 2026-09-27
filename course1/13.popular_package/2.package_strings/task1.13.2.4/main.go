package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func generateActivationKey() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	keyParts := [19]string{}
	symbols := strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")
	for i := 0; i < 19; i++ {
		keyParts[i] = symbols[rand.Intn(len(symbols))]
		if i != 0 && i%5 == 0 {
			keyParts[i-1] = "-"
		}
	}

	return strings.Join(keyParts[:], "")
}

func main() {
	activationKey := generateActivationKey()
	fmt.Println(activationKey) // UQNI-NYSI-ZVYB-ZEFQ
}
