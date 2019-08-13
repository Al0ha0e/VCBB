from bottle import request, response, route, run, template
import time
import json
import redis

red = redis.Redis(host='localhost', port=6379)


def execute(keys, partitionCnt, code):
    ret = []
    l = len(keys)
    seg = l//partitionCnt
    i = 0
    print(code)
    while i < partitionCnt:
        args = {'input': red.mget(*keys[i*seg:(i+1)*seg]), 'output': []}
        exec(code, args)
        print(args['output'])
        ret.append(args['output'])
        i += 1
    return ret


@route('/hello/<name>', method="post")
def index(name):
    print("DDDDDB", request.json["keys"],
          request.json["partitionCnt"], request.json["code"])
    res = execute(request.json["keys"],
                  request.json["partitionCnt"], request.json["code"])
    print(res)
    response.content_type = "application/json"
    return json.dumps(res)
    # return template('<b>Hello {{name}}</b>!', name=name)


run(host='localhost', port=8080)
