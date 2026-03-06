package wait

import (
	IO "github.com/IBM/fp-go/v2/io"
)

// logState returns a logging IO action for use with IOE.ChainFirstIOK.
// Reads Logger and State directly from PollContext — no partial application needed.
func logState(ctx PollContext) IO.IO[struct{}] {
	return func() struct{} {
		ctx.Logger.Info("polling",
			"job", ctx.Name,
			"state", string(ctx.State),
		)
		return struct{}{}
	}
}
