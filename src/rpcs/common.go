package rpcs

type RequestBody struct {
	Method string			`json:method`
	Params []interface{}	`json:params`
	Id string				`json:id`
}
