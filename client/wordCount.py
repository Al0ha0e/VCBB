import preprocessor
import client
import json
import re

data = []
i = 0
while i < 3:
    f = open("wc"+str(i)+".txt", "r", encoding="utf-8")
    s = f.read()
    data.append(s.encode('utf-8'))
    f.close()
    i += 1

cpns = []


def imp1(id):
    return [id]


def omp1(id):
    return [id+3]


code1 = """
import re
import json
s=input[0].decode('utf-8')
fil = re.compile("[^a-z|^A-Z]")
text = fil.sub(" ", s.strip())
words = text.split(' ')
count = {}
for word in words:
    if not word.isalpha():
        continue
    if word in count:
        count[word] += 1
    else:
        count[word] = 1
output = [json.dumps(count).encode('utf-8')]
"""

cpns.append(preprocessor.computeNode(3, 1, 1, imp1, omp1, code1))


def imp2(id):
    return [3, 4]


def omp2(id):
    return [6]


code2 = """
import json
count1 = json.loads(input[0].decode('utf-8'))
count2 = json.loads(input[1].decode('utf-8'))
for k,v in count1.items():
    if k in count2:
        count2[k] += v
    else:
        count2[k] = v
print(count2)
output = [json.dumps(count2).encode('utf-8')]
"""

cpns.append(preprocessor.computeNode(1, 2, 1, imp2, omp2, code2))


def imp3(id):
    return[5, 6]


def omp3(id):
    return[7]


code3 = code2

cpns.append(preprocessor.computeNode(1, 2, 1, imp3, omp3, code3))

client.commitScheduleGraph(8, cpns, data)
