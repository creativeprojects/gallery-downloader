package download

import (
	"net/url"
	"path"
	"strings"
)

func joinURL(first, second *url.URL) *url.URL {
	join := &url.URL{
		Scheme: first.Scheme,
		Host:   first.Host,
	}
	if strings.HasPrefix(second.Path, "/") {
		// This is an absolute path from the root
		join.Path = second.Path
		return join
	}
	if !strings.HasSuffix(first.Path, "/") {
		join.Path = path.Dir(first.Path)
		if join.Path == "." {
			join.Path = "/"
		}
		join.Path = path.Join(join.Path, second.Path)
		return join
	}
	join.Path = path.Join(first.Path, second.Path)
	return join
}
