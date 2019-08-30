from bottle import request, response, route, run, template
import time
import json as JSON
import redis
import hashlib

red = redis.Redis(host='localhost', port=6379)

# TODO: RUNTIME ERROR REPORT


def execute(offset, keys, partitionCnt, code):
    ret = []
    l = len(keys)
    i = 0
    # print(code)
    while i < l:
        args = {'partitionId': offset+i,
                'input': red.mget(*keys[i]), 'output': []}
        exec(code, args)
        ans = args['output']
        #print(i, ans)
        #print("BEFORE", ans, args['output'])
        j = 0
        for obj in ans:
            sha3 = hashlib.sha3_256()
            sha3.update(obj)
            key = sha3.hexdigest()
            red.set(key, obj)
            ans[j] = key
            j += 1
        #print(i, "AFTER", ans)
        ret.append(ans)
        i += 1
    return ret


@route('/execute', method="post")
def executer():
    #print("DDDDDB", request.json["keys"],request.json["partitionCnt"], request.json["code"])
    res = execute(request.json["partitionIdOffset"], request.json["keys"],
                  request.json["partitionCnt"], request.json["code"])
    # print(res)
    response.content_type = "application/json"
    return JSON.dumps(res)
    # return template('<b>Hello {{name}}</b>!', name=name)


run(host='localhost', port=8080)
