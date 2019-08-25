import preprocessor
import client
import json

cnps = []


def imp1(id):
    return [0]


def omp1(id):
    return [1, 2]


code1 = """
t = input[0]
output=[t+b'1919810',t+b'114514']
"""

cnps.append(preprocessor.computeNode(1, 1, 2, imp1, omp1, code1))


def imp2(id):
    return [2]


def omp2(id):
    return [3, 4]


code2 = """
t = input[0]
output=[t+b'1919810',t+b'114514']
"""

cnps.append(preprocessor.computeNode(1, 1, 2, imp2, omp2, code2))


def imp3(id):
    return [1, 3]


def omp3(id):
    return [5]


code3 = """
output=[input[0]+input[1]]
"""

cnps.append(preprocessor.computeNode(1, 2, 1, imp3, omp3, code3))


def imp4(id):
    return [5, 4]


def omp4(id):
    return [6]


code4 = """
output=[input[0]+input[1]]
"""

cnps.append(preprocessor.computeNode(1, 2, 1, imp4, omp4, code4))

data = [str(114514).encode('utf-8')]

client.commitScheduleGraph(7, cnps, data)
