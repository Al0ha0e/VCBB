import preprocessor
import redis
import requests
import hashlib
import json

red = redis.Redis(db=1)


def commitScheduleGraph(dataNodeCnt, computeNodes, oriData):
    red = redis.Redis(db=1)
    oriDataMap, schGraph = preprocessor.preprocess(dataNodeCnt, computeNodes)
    oriDataHash = {}
    for key in oriDataMap.keys():
        data = oriData[key]
        sha3 = hashlib.sha3_256()
        sha3.update(data)
        hs = sha3.hexdigest()
        oriDataHash[oriDataMap[key]] = hs
        red.set(hs, data)
    jsonSchGraph = []
    for node in schGraph:
        jsonSchGraph.append(node.toDict())
    req = json.dumps({'oriDataHash': oriDataHash, 'schGraph': jsonSchGraph})
    requests.post("http://localhost:8080/commitSchGraph", data=req)
