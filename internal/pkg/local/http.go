package local

import (
	"fmt"
	"net/http"
)

func isResponding(cmp Component) string {
	if cmp.testUrl == "" {
		return ""
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(cmp.testUrl)
	if err != nil {
		// handle error
		fmt.Println(err)
		return "ERR"
	}
	defer resp.Body.Close()
	return fmt.Sprintf("%v", resp.StatusCode)
}
