package masks

import (
	"fmt"
	"regexp"
)

type Type string

const (
	Inclusive Type = "inclusive"
	Exclusive Type = "exclusive"
)

type Pattern struct {
	Type  Type
	Regex string
	re    *regexp.Regexp
}

func (p *Pattern) Match(fileName string) (bool, error) {
	if p.re == nil {
		var err error
		if p.re, err = regexp.Compile(p.Regex); err != nil {
			return false, fmt.Errorf("invalid regex for %s pattern: %q", p.Type, p.Regex)
		}
	}
	return p.re.Match([]byte(fileName)), nil
}
