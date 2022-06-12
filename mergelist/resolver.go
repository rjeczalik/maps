package mergelist

import "net/url"

type Resolver interface {
	Resolve(*url.URL) (Iter, error)
}
