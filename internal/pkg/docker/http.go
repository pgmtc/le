package docker

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"net/http"
	"strconv"
	"time"
)

func isResponding(cmp common.Component) (result string, resultErr error) {
	timeout := time.Duration(3 * time.Second)
	if cmp.TestUrl == "" {
		result = ""
		return
	}

	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(cmp.TestUrl)
	if err != nil {
		result = "ERR"
		resultErr = err
		return
	}
	defer resp.Body.Close()
	result = strconv.Itoa(resp.StatusCode)
	return
}
