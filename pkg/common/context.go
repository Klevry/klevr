package common

var (
	BaseContext = NewContext()
)

type Context struct {
	variables map[string]interface{}
}

func NewContext() *Context {
	return &Context{
		make(map[string]interface{}),
	}
}

func FromContext(ctx *Context) *Context {
	tMap := make(map[string]interface{})

	for k, v := range ctx.variables {
		tMap[k] = v
	}

	return &Context{
		tMap,
	}
}

func (c *Context) Put(key string, value interface{}) {
	c.variables[key] = value
}

func (c *Context) Get(key string) interface{} {
	return c.variables[key]
}
