package local

import (
	"fmt"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"net/http"
	"time"
)

func isResponding(cmp common.Component) string {
	timeout := time.Duration(3 * time.Second)
	if cmp.TestUrl == "" {
		return ""
	}

	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(cmp.TestUrl)
	if err != nil {
		// handle error
		fmt.Println(err)
		return "ERR"
	}
	defer resp.Body.Close()
	return fmt.Sprintf("%v", resp.StatusCode)
}
