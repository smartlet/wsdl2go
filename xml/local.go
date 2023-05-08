package xml

import "strings"

func local(name string) string {
	idx := strings.IndexByte(name, ':')
	if idx == -1 {
		return name
	}
	return name[idx+1:]
}
