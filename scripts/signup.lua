math.randomseed(os.time())

request = function()
    local username = "user" .. math.random(100000, 999999)
    local body = string.format('{"name":"%s","password":"hello$1234","confirmPassword":"hello$1234"}', username)
    wrk.body = body
    wrk.headers["Content-Type"] = "application/json"
    return wrk.format("POST", "/user/signup")
end