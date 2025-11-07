package pagination

// SliceParams holds parameters for slicing results
type SliceParams struct {
	Cursor int64
	Limit  int64
	Total  int64
}

// SliceResult holds the result of slicing operation
type SliceResult struct {
	Start      int64
	End        int64
	NextCursor int64
}

// CalculateSlice calculates pagination slice boundaries and next cursor
func CalculateSlice(params SliceParams) SliceResult {
	start := params.Cursor
	if start < 0 {
		start = 0
	}
	if start > params.Total {
		start = params.Total
	}

	end := start + params.Limit
	if end > params.Total {
		end = params.Total
	}

	// Calculate next cursor
	nextCursor := int64(0)
	if end < params.Total {
		nextCursor = end
	}

	return SliceResult{
		Start:      start,
		End:        end,
		NextCursor: nextCursor,
	}
}
