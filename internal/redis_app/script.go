package redis_app

import "github.com/go-redis/redis/v8"

var _writeScript = redis.NewScript(`
	local num_arg = #ARGV/2
	local field_index = 1
	local value_index = 1
	for i=1, num_arg do
		field_index = ((i-1)*2) + 1
		value_index = field_index + 1
		redis.call("HSET", KEYS[i], ARGV[field_index], ARGV[value_index])
	end

	return 0
`)

var _writeScript2 = redis.NewScript(`
	local num_arg = #ARGV/3
	local key_index = 1
	local field_index = 1
	local value_index = 1
	for i=1, num_arg do
		key_index = ((i-1)*3) + 1
		field_index = key_index + 1
		value_index = field_index + 1
		redis.call("HSET", ARGV[key_index], ARGV[field_index], ARGV[value_index])
	end

	return 0
`)

var _readScript = redis.NewScript(`
	local num_arg = #ARGV

	local retval = {}
	for i=1, num_arg do
		local doc = redis.call("HGETALL", ARGV[i])
		if #doc > 0 then
			table.insert(retval, ARGV[i])
			table.insert(retval, doc)
		end
	end

	if #retval == 0 then
		return -1
	end

	return retval
`)

var _counterScript = redis.NewScript(`
	local doc_path = KEYS[1]
	local parent_id = ARGV[1]
	local created_at = ARGV[2]
	local key = ARGV[3]
	local val = ARGV[4]

	local keyExist = redis.call("EXISTS", doc_path)
	if keyExist == 0 then
		redis.call("HMSET", doc_path, "$rev", 1, "$createdAt", created_at, "$parentId", parent_id, key, val)
	else
		redis.call("HINCRBY", doc_path, "$rev", 1)
		redis.call("HINCRBYFLOAT", doc_path, key, val)
	end

	local snapshot = redis.call("HGETALL", doc_path)
	if #snapshot == 0 then
		return -1
	end
	
	local retval = {}
	table.insert(retval, doc_path)
	table.insert(retval, snapshot)
	return retval
`)
