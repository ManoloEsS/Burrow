package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ManoloEsS/burrow_prototype/internal/config"
	"github.com/ManoloEsS/burrow_prototype/internal/models"
)

func FormatInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanned := scanner.Scan()
	if !scanned {
		return ""
	}

	line := scanner.Text()
	line = strings.TrimSpace(line)
	return line
}

func InputToReq(cfg *config.Config, request *models.Request) {
	fmt.Println("Input a method (upper or lowercase)")
	fmt.Print("> ")
	request.ParseMethod(FormatInput())

	fmt.Println("Input a url (default localhost:8080)")
	fmt.Print("> ")
	request.ParseUrl(cfg, FormatInput())

	fmt.Println("Input headers (key:value key:value)")
	fmt.Print("> ")
	request.ParseHeaders(FormatInput())

	fmt.Println("Input params (key:value key:value)")
	fmt.Print("> ")
	request.ParseParams(FormatInput())

	fmt.Println("Body (string)")
	fmt.Print("> ")
	request.ParseBody(FormatInput())

	fmt.Println("Auth (placeholder string)")
	fmt.Print("> ")
	request.ParseAuth(FormatInput())
}
