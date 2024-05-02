"""
meta:prev_req: Some previous request
doc:url: http://foobar.com
doc:method: POST
"""
url = "http://foobar.com"
headers = { "X-Foo": "Bar", "X-Foos": [ "Bar1", "Bar2" ] }
method = "POST"
body = {
    "id": 1,
    "amount": 1.2001,
    "name": "Jane"
}

