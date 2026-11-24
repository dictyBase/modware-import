// Package cli registers the phenoannotation commands.
package cli

import (
	"context"

	cli "github.com/urfave/cli/v3"
)

//go:generate go run github.com/IBM/fp-go/v2@v2.13.2 lens --dir . --filename gen_lens.go

// arguments is the boundary record of one CLI invocation.
// fp-go:Lens
type arguments struct {
	Ctx context.Context
	Cmd *cli.Command
}
