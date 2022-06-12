package mergelist

var defaultBuilder = newBuilder()

type Builder struct {
}

func newBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Resolve(l List) Iter {
	return nil
}
