package ontology

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dictyBase/go-obograph/graph"
	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
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

type updateTermRowProperties struct {
	*termRowProperties
	RowId int32
}

type exisTermRowResp struct {
	Exist        bool
	IsDeprecated bool
	RowId        int32
}

type ontologyRow struct {
	Id         int32  `json:"id"`
	TermId     string `json:"term_id"`
	Name       string `json:"name"`
	IsObsolete bool   `json:"is_obsolete"`
}

type ontologyListRows struct {
	Count    int32                 `json:"count"`
	Next     client.NullableString `json:"next"`
	Previous client.NullableString `json:"previous"`
	Results  []*ontologyRow        `json:"results"`
}

func LoadNewOrUpdate(args *LoadProperties) error {
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
		existResp, err := existTermRow(&termRowProperties{
			Term:    term,
			Host:    args.Client.GetConfig().Host,
			Token:   args.Token,
			TableId: args.TableId,
		})
		if err != nil {
			return err
		}
		if existResp.Exist {
			if existResp.IsDeprecated == term.IsDeprecated() {
				args.Logger.Debugf("term %s has no change", string(term.ID()))
				continue
			}
			err = updateTermRow(&updateTermRowProperties{
				RowId: existResp.RowId,
				termRowProperties: &termRowProperties{
					Term:    term,
					Host:    args.Client.GetConfig().Host,
					Token:   args.Token,
					TableId: args.TableId,
				},
			})
			if err != nil {
				return err
			}
			args.Logger.Infof("updated row with term %s", string(term.ID()))
			continue
		}
		err = addTermRow(&termRowProperties{
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
		err := addTermRow(&termRowProperties{
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

func updateTermRow(args *updateTermRowProperties) error {
	payload := map[string]interface{}{
		"Is_obsolete": termStatus(args.Term),
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
	httpapi.CommonHeader(req, args.Token, "Token")
	res, err := httpapi.ReqToResponse(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func existTermRow(args *termRowProperties) (*exisTermRowResp, error) {
	term := string(args.Term.ID())
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/api/database/rows/table/%d/?user_field_names=true&size=1&search=%s",
			args.Host,
			args.TableId,
			term,
		), nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error in creating requst %s", err)
	}
	httpapi.CommonHeader(req, args.Token, "Token")
	res, err := httpapi.ReqToResponse(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	rowsResp := &ontologyListRows{}
	if err := json.NewDecoder(res.Body).Decode(rowsResp); err != nil {
		return nil, fmt.Errorf("error in decoding json response %s", err)
	}
	existResp := &exisTermRowResp{Exist: false}
	if rowsResp.Count > 0 {
		existResp.Exist = true
		existResp.IsDeprecated = rowsResp.Results[0].IsObsolete
		existResp.RowId = rowsResp.Results[0].Id
	}
	return existResp, nil
}

func addTermRow(args *termRowProperties) error {
	term := args.Term
	payload := map[string]interface{}{
		"Id":          string(term.ID()),
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
	httpapi.CommonHeader(req, args.Token, "Token")
	res, err := httpapi.ReqToResponse(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func termStatus(term graph.Term) string {
	if term.IsDeprecated() {
		return "true"
	}
	return "false"
}
