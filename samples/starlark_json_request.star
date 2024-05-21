"""
prev_req: ""
doc:url: https://httpbin.org/anything
doc:method: POST
"""
url = "https://httpbin.org/anything"
method = "POST"
headers = { "X-Foo": "bar", "Content-Type": "application/json" }
body = { "id": 1474, "bar": [
    {"name": "Joe"},
    {"name": "Jane"},
] }
