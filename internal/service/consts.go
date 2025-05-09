package service

const CONTROLLER_ACTION_TMPL = `
// @http_method %s
// @middleware
func (%s) %s(c *gin.Context) {
	//TODO: edit
}
`
