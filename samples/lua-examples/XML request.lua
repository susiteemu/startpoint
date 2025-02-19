return {
	url = "http://httpbin.org/anything",
	method = "POST",
	headers = {
		["Content-Type"] = "application/xml",
	},
	body = [[
    <root>
      <id>1</id>
      <name>Jane</name>
    </root>
  ]],
}

