local setName = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local clearBefore = now - window
redis.call('ZREMRANGEBYSCORE', setName, 0, clearBefore)
local amount = redis.call('ZCARD', setName)
if amount < limit then
	redis.call('ZADD', setName, now, now)
end
redis.call('EXPIRE', setName, window)
return limit - amount

