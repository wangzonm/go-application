wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

-- 初始化计数器
counter = 0

-- 定义协程函数用于生成计数
function count_generator()
    while true do
        counter = counter + 1
        coroutine.yield(counter)
    end
end

-- 创建一个新的协程
count_co = coroutine.create(count_generator)

-- 定义请求处理函数
function request()
    local timestamp = os.time()
    local random_value = math.random(1000)

    -- 使用协程获取当前计数
    local _, order_id = coroutine.resume(count_co)
    order_id = "BTC-USDT-SELL" .. string.format("%04d", order_id)

    local body = [[
        {
            "action": "create",
            "symbol": "BTC-USDT",
            "order_id": "]] .. order_id .. [[",
            "side": 2,
            "type": 8,
            "amount": ]] .. random_value .. [[,
            "price": ]] .. random_value .. [[,
            "timestamp": ]] .. timestamp .. [[
        }
    ]]

    return wrk.format("POST", nil, nil, body)
end

-- 请求完成时运行
function done(summary, latency, requests)
    -- 输出测试结果统计信息
    io.write("Total requests: ", summary.requests, "\n")
    io.write("Failed requests: ", summary.errors.status, "\n")
    io.write("Request latency (ms):\n")
    for _, p in pairs({ 50, 90, 99, 99.9 }) do
        n = latency:percentile(p)
        io.write(string.format("%g%%,%d\n", p, n))
    end
end
