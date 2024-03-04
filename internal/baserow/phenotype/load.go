package phenotype

import (
	"github.com/sirupsen/logrus"
)

type LoaderProperties struct {
	File    string
	TableId int
	Token   string
	Host    string
	Logger  *logrus.Entry
}

type PhenotypeLoader struct {
	file    string
	tableId int
	token   string
	host    string
	logger  *logrus.Entry
}

func NewPhenotypeLoader(args LoaderProperties) *PhenotypeLoader {
	return &PhenotypeLoader{
		file:    args.File,
		tableId: args.TableId,
		token:   args.Token,
		host:    args.Host,
		logger:  args.Logger,
	}
}

type ParseNameToDateFeedback struct {
	Err error
}

func (loader *PhenotypeLoader) Load(args *LoaderProperties) error {
	return nil
}
