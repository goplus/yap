get "/", ctx => {
	ctx.html `<html><body>Hello, YAP!</body></html>`
}
get "/p/:id", ctx => {
	ctx.json {
		"id": ctx.param("id"),
	}
}

run "localhost:8080"
