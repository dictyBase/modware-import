package regexp

import "regexp"

var DateRegxp = regexp.MustCompile(`^(\w{2})-(\w{3})-(\w{2})`)
