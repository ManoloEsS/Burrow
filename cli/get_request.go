package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type HttpRequest struct {
	Headers string
}

type RequestMsg struct {
	Body    string
	Headers string
	Err     error
}

func GetRequest(url string) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return RequestMsg{Err: err}
		}

		res, err := client.Do(req)
		if err != nil {
			return RequestMsg{Err: err}
		}
		defer res.Body.Close()

		bodyByte, err := io.ReadAll(res.Body)
		if err != nil {
			return RequestMsg{Err: err}
		}

		return RequestMsg{
			Body: string(bodyByte),
		}
	}
}

func GetRequestPrint(name string) {
	fmt.Printf("this is the request %s\n", name)
}
