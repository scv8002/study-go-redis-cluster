package redis_app

import "github.com/go-redis/redis/v8"

var _writeScript = redis.NewScript(`
	local num_arg = #ARGV/2
	local field_index = 1
	local value_index = 1
	for i=1, num_arg do
		field_index = ((i-1)*2) + 1
		value_index = field_index + 1
		redis.call("HMSET", KEYS[i], ARGV[field_index], ARGV[value_index])
	end

	return 0
`)
