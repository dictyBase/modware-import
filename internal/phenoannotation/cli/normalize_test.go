package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	cli "github.com/urfave/cli/v3"

	norm "github.com/dictyBase/modware-import/internal/phenoannotation/normalize"
)

func TestCancelledRun(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cmd := &cli.Command{
		Name:   "normalize",
		Action: RunNormalize,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "folder"},
			&cli.StringFlag{Name: "out"},
			&cli.StringFlag{Name: "sheet"},
		},
	}
	err := cmd.Run(ctx, []string{"normalize", "--folder", dir, "--out", out})
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if _, err := os.Stat(out); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}

func TestExitCode(t *testing.T) {
	cases := map[error]int{
		nil:                                  exitOK,
		ErrReview:                            exitReview,
		fmt.Errorf("wrapped: %w", ErrReview): exitReview,
		io.ErrUnexpectedEOF:                  exitFailed,
	}
	for err, want := range cases {
		if got := ExitCode(err); got != want {
			t.Fatalf("%v: got %d want %d", err, got, want)
		}
	}
}

func TestReviewStatus(t *testing.T) {
	var s norm.State
	if E.IsLeft(requireNoRejects(s)) {
		t.Fatal("clean batch rejected")
	}
	_, err := E.UnwrapError(requireNoRejects(stateRejected.Set(1)(s)))
	if !errors.Is(err, ErrReview) {
		t.Fatal(err)
	}
}
