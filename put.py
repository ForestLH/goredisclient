import redis

def put():
	con_pool = redis.ConnectionPool(host='127.0.0.1', port=6379, decode_responses=True)
	handle = redis.StrictRedis(connection_pool=con_pool, charset='utf-8')
	
	key = 'test'
	
	batch = []
	for i in range(1000000):
		batch.append(str(i))
		
		if len(batch)>10000:
			handle.rpush(key, *batch)
			batch = []
			
	if len(batch)>0:
		handle.rpush(key, *batch)
		
if __name__== '__main__':
	put()
