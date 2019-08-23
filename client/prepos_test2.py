import preprocessor
import client
import json

cnps = []


def imp1(id):
    return [id*3, id*3+1, id*3+2]


def omp1(id):
    return [id*4+9, id*4+10, id*4+11, id*4+12]


code1 = """
"""

cnps.append(preprocessor.computeNode(2, 3, 4, imp1, omp1, code1))
