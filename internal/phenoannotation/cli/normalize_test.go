package cli

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	E "github.com/IBM/fp-go/v2/either"

	norm "github.com/dictyBase/modware-import/internal/phenoannotation/normalize"
)

func TestCancelledRun(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := Command(io.Discard).Run(ctx, []string{
		"phenoannotation", "normalize", "--folder", dir, "--out", out,
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if _, err := os.Stat(out); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}

func TestNeedsSubcommand(t *testing.T) {
	err := Command(io.Discard).Run(context.Background(), []string{"phenoannotation"})
	if !errors.Is(err, errNoSub) {
		t.Fatal(err)
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
