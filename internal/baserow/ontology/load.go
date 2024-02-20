import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/dictyBase/go-obograph/graph"
	"github.com/dictyBase/modware-import/internal/baserow/client"
)
func addTermRow(args *addTermRowProperties) error {
	term := args.Term
	payload := map[string]interface{}{
		"Id":          term.ID(),
		"Name":        term.Label(),
		"Is_obsolete": termStatus(term),
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error in encoding body %s", err)
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"%s/api/database/rows/table/%d/?user_field_names=true",
			args.Host,
			args.TableId,
		),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error in creating request %s ", err)
	}
	commonHeader(req, args.Token)
	res, err := reqToResponse(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

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

func termStatus(term graph.Term) string {
	if term.IsDeprecated() {
		return "true"
	}
	return "false"
}
