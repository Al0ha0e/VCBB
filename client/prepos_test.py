import preprocessor
import json

cpns = []


def imp1(id):
    return [id*2, id*2+1]


def omp1(id):
    return [id*3+10, id*3+11, id*3+12]


cpns.append(preprocessor.computeNode(5, 2, 3, imp1, omp1, "ps1"))


def imp2(id):
    i = 0
    ret = []
    while i < 5:
        ret.append(10+id+3*i)
        i += 1
    return ret


def omp2(id):
    return [25+id]


cpns.append(preprocessor.computeNode(3, 5, 1, imp2, omp2, "ps2"))
ori, schndoes = preprocessor.preprocess(28, cpns)
print(ori)
for sn in schndoes:
    sn.show()
jsonschs = []
for node in schndoes:
    jsonschs.append(node.toDict())
print(json.dumps(jsonschs))
