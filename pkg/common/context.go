package common

var store = make(map[string]interface{})

// ContextPut put value with key
func ContextPut(key string, value interface{}) {
	store[key] = value
}

// ContextGet get value with key
func ContextGet(key string) interface{} {
	return store[key]
}

// ContextGetString get value as string with key
func ContextGetString(key string) string {
	return (store[key]).(string)
}
