--[[
prev_req: Fetch token
--]]
local auth
auth = "Bearer " .. prevResponse.body["access_token"]
return {
	url = "http://localhost:8000/auth/oauth2/users/me/",
	method = "GET",
	headers = { Authorization = auth },
	options = { printRequest = true },
}
