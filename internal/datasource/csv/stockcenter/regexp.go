package stockcenter

import "regexp"

var dateRegxp = regexp.MustCompile(`^(\w{2})-(\w{3})-(\w{2})`)
