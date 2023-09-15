-- 设置请求方法和头部
wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"

-- 定义请求处理函数
function request()
    local uid = math.random(1, 1000)  -- 随机生成uid
    local url = "/v1/activity/mutil_without_pipline?uid=" .. uid

    return wrk.format("GET", url)
end
