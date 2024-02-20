import (
	"fmt"
	"io"
	"net/http"
)
func reqToResponse(creq *http.Request) (*http.Response, error) {
	client := &http.Client{}
	uresp, err := client.Do(creq)
	if err != nil {
		return uresp, fmt.Errorf("error in making request %s", err)
	}
	if uresp.StatusCode != 200 {
		cnt, err := io.ReadAll(uresp.Body)
		if err != nil {
			return uresp, fmt.Errorf(
				"error in response and the reading the body %d %s",
				uresp.StatusCode,
				err,
			)
		}
		return uresp, fmt.Errorf(
			"unexpected error response %d %s",
			uresp.StatusCode,
			string(cnt),
		)
	}
	return uresp, nil
}

func commonHeader(lreq *http.Request, token string) {
	lreq.Header.Set("Content-Type", "application/json")
	lreq.Header.Set("Accept", "application/json")
	lreq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}
