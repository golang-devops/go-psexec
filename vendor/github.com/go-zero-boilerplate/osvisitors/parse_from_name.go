package osvisitors

import (
	"fmt"
	"strings"
)

func ParseFromName(name string) (OsType, error) {
	for _, os := range AllList {
		visitor := &GoOSNameVisitor{}
		os.Accept(visitor)
		if strings.EqualFold(name, visitor.Name) {
			return os, nil
		}
	}
	return nil, fmt.Errorf("github.com/go-zero-boilerplate/osvisitors does not currently support OS name '%s'", name)
}
