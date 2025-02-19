return {
	-- Request url
	url = "http://localhost:8000/multipart-form",
	-- HTTP method
	method = "POST",
	headers = {
		["Content-Type"] = "multipart/form-data",
	},
	body = {
		title = "File title",
		file = "@luatests/test.txt",
	},
}
