package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type WorkspaceManager struct {
	Logger *logrus.Entry
	Token  string
	Host   string
}

func (wkm *WorkspaceManager) ListWorkspaceURL() string {
	return fmt.Sprintf("https://%s/api/workspaces/", wkm.Host)
}
