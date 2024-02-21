package ontology

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dictyBase/go-obograph/graph"
	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/sirupsen/logrus"
)

type LoadProperties struct {
	File    string
	TableId int
	Token   string
	Client  *client.APIClient
	Logger  *logrus.Entry
}

type termRowProperties struct {
	Term    graph.Term
	Host    string
	Token   string
	TableId int
}

}

func LoadNew(args *LoadProperties) error {
	rdr, err := os.Open(args.File)
	if err != nil {
		return fmt.Errorf("error in opening file %s %s", args.File, err)
	}
	defer rdr.Close()
	grph, err := graph.BuildGraph(rdr)
	if err != nil {
		return fmt.Errorf(
			"error in building graph from file %s %s",
			args.File,
			err,
		)
	}
	for _, term := range grph.Terms() {
		err := addTermRow(&addTermRowProperties{
			Term:    term,
			Host:    args.Client.GetConfig().Host,
			Token:   args.Token,
			TableId: args.TableId,
		})
		if err != nil {
			return err
		}
		args.Logger.Infof("add row with id %s", term.ID())
	}
	return nil
}

func existTermRow(args *existTermRowProperties) (bool, error) {
	ok := false
	term := args.Term
	rows, resp, err := args.Client.DatabaseTableRowsApi.
		ListDatabaseTableRows(args.Ctx, int32(args.TableId)).
		Size(1).
		UserFieldNames(true).
		Search(string(term.ID())).
		Execute()
	if err != nil {
		return ok, fmt.Errorf(
			"error in checking presence of term %s %s",
			string(term.ID()),
			err,
		)
	}
	defer resp.Body.Close()
	if rows.Count > 0 {
		ok = true
func updateTermRow(args *updateTermRowProperties) error {
	term := args.Term
	payload := map[string]interface{}{
		"Name":        term.Label(),
		"Is_obsolete": termStatus(term),
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error in encoding body %s", err)
	}
	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf(
			"%s/api/database/rows/table/%d/%d/?user_field_names=true",
			args.Host,
			args.TableId,
			args.RowId,
		),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error in creating requst %s", err)
	}
	commonHeader(req, args.Token)
	res, err := reqToResponse(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func existTermRow(args *termRowProperties) (*exisTermRowResp, error) {
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
