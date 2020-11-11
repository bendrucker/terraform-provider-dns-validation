package provider

import "context"

type Resolver interface {
	LookupHost(context.Context, string) ([]string, error)
}
