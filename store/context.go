package store

import "golang.org/x/net/context"

const Key = "store"

type Setter interface {
	Set(string, interface{})
}

func FromContext(c context.Context) Store {
	return c.Value(Key).(Store)
}

func ToContext(c Setter, store Store) {
	c.Set(Key, store)
}
