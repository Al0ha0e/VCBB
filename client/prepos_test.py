import preprocessor
import client
import json

cpns = []


def imp1(id):
    return [id*2, id*2+1]


def omp1(id):
    return [id*3+10, id*3+11, id*3+12]


code1 = """
print(input)
tmp = input[0]+input[1]
print(tmp)
output = [tmp+b's',tmp+b'z',tmp+b'h']
"""

cpns.append(preprocessor.computeNode(5, 2, 3, imp1, omp1, code1))


def imp2(id):
    i = 0
    ret = []
    while i < 5:
        ret.append(10+id+3*i)
        i += 1
    return ret


def omp2(id):
    return [25+id]


code2 = """
output = [input[0]+input[1]+input[2]+input[3]+input[4]]
"""

cpns.append(preprocessor.computeNode(3, 5, 1, imp2, omp2, code2))
# ori, schndoes = preprocessor.preprocess(28, cpns)
# print(ori)
# for sn in schndoes:
#     sn.show()
# jsonschs = []
# for node in schndoes:
#     jsonschs.append(node.toDict())
# print(json.dumps(jsonschs))
data = []
i = 0
while i < 10:
    data.append(str(i).encode('utf-8'))
    i += 1

client.commitScheduleGraph(28, cpns, data)
