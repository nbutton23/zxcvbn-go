package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nbutton23/zxcvbn-go"
)

func main() {

	fmt.Println("Enter password:")
	reader := bufio.NewReader(os.Stdin)
	password, _ := reader.ReadString('\n')
	//password := "Testaaatyhg890l33t"

	passwordStenght := zxcvbn.PasswordStrength(password, nil)

	fmt.Printf(
		`Password score    (0-4): %d
Estimated entropy (bit): %f
Estimated time to crack: %s%s`,
		passwordStenght.Score,
		passwordStenght.Entropy,
		passwordStenght.CrackTimeDisplay, "\n",
	)
}
