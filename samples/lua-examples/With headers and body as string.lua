return {
	url = "https://httpbin.org/anything",
	method = "POST",
	headers = {
		Accept = "application/json",
		["Content-Type"] = "application/json",
		["X-Custom-Header"] = "Some custom value",
	},
	body = '{"id": 1, "name": "Jane"}',
}

