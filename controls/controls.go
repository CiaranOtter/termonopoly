package controls

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func YesNo(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (Y/n): ", question)

	input, err := reader.ReadString('\n')

	if err != nil {
		return YesNo(question)
	}

	if strings.Contains(strings.ToLower(input), "y") {
		return true
	}

	return false
}
