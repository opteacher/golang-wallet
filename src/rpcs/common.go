package rpcs

var PKG struct {}

type RequestBody struct {
	Method string			`json:method`
	Params []interface{}	`json:params`
	Id string				`json:id`
}
