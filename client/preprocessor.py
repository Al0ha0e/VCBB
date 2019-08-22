import random
import string


class dataNode:
    def __init__(self, id):
        self.id = id
        self.input = "ori"
        self.output = []


class computeNode:
    def __init__(self, partitionCnt, inputCnt, outputCnt, inputMapper, outputMapper, processor):
        self.partitionCnt = partitionCnt
        self.inputCnt = inputCnt
        self.outputCnt = outputCnt
        self.inputMapper = inputMapper
        self.outputMapper = outputMapper
        self.processor = processor


class schNode:
    def __init__(self, id, code, ptCnt, ptIDOffset, minAnswerCount,  cpNode, baseTest="", hdreq=""):
        self.id = id
        self.code = code
        self.baseTest = baseTest
        self.hardwareRequirement = hdreq
        self.partitionCnt = ptCnt
        self.partitionIDOffset = ptIDOffset
        self.dependencies = {}
        self.inputCnt = cpNode.inputCnt
        self.inputMap = {}
        self.output = []
        self.indeg = 0
        self.outdeg = 0
        #self.inNodes = []
        self.outNodes = []
        self.minAnswerCount = minAnswerCount
        self.cpNode = cpNode

    def toDict(self):
        ret = {'id': self.id,
               'code': self.code,
               'baseTest': self.baseTest,
               'hardwareRequirement': self.hardwareRequirement,
               'partitionCnt': self.partitionCnt,
               'partitionIDOffset': self.partitionIDOffset,
               'dependencies': self.dependencies,
               'inputCnt': self.inputCnt,
               'inputMap': self.inputMap,
               'output': self.output,
               'indeg': self.indeg,
               'outdeg': self.outdeg,
               'outNodes': self.outNodes,
               'minAnswerCount': self.minAnswerCount}
        return ret

    def show(self):
        print('------------------------------')
        print(self.id)
        print(self.code, self.baseTest, self.hardwareRequirement)
        print(self.partitionCnt, self.partitionIDOffset, self.minAnswerCount)
        print(self.inputMap)
        print(self.output)
        print(self.dependencies, self.indeg)
        print(self.outNodes, self.outdeg)


def randomId(l):
    s = set()
    while True:
        while True:
            r = ''.join(random.sample(string.ascii_letters+string.digits, l))
            if r not in s:
                s.add(r)
                yield r


def preprocess(dataNodeCnt, computeNodes):
    dataNodes = []
    dataNodesRev = {}
    schNodes = []
    i = 0
    rid = randomId(32)
    while i < dataNodeCnt:
        tid = next(rid)
        dataNodes.append(dataNode(tid))
        dataNodesRev[tid] = i
        i += 1
    nid = randomId(32)
    for node in computeNodes:
        cnt = node.partitionCnt
        offset = 0
        while cnt > 0:
            num = 3
            if cnt < 3:
                num = cnt
            cnt -= 3
            schNodes.append(
                schNode(next(nid), node.processor, num, offset, 1, node))
            offset += 3
    for node in schNodes:
        inmp = node.cpNode.inputMapper
        outmp = node.cpNode.outputMapper
        st = node.partitionIDOffset
        i = 0
        while i < node.partitionCnt:
            oriInpt = inmp(st+i)
            j = 0
            while j < len(oriInpt):
                dn = dataNodes[oriInpt[j]]
                node.inputMap[dn.id] = {'x': i, 'y': j}
                dn.output.append(node.id)
                j += 1
            oriOutput = outmp(st+i)
            j = 0
            output = []
            while j < len(oriOutput):
                dn = dataNodes[oriOutput[j]]
                output.append(dn.id)
                dn.input = node.id
                j += 1
            node.output.append(output)
            i += 1

    for node in schNodes:
        for key in node.inputMap:
            dn = dataNodes[dataNodesRev[key]]
            if dn.input not in node.dependencies:
                node.dependencies[dn.input] = []
            node.dependencies[dn.input].append(key)
        ks = node.dependencies.keys()
        node.indeg = len(list(ks))
        if node.indeg == 1 and ("ori" in ks):
            node.indeg = 0
        opst = set()
        for op in node.output:
            for opid in op:
                dn = dataNodes[dataNodesRev[opid]]
                for ipid in dn.output:
                    opst.add(ipid)
        i = 0
        for ipid in opst:
            node.outNodes.append(ipid)
            i += 1
        node.outdeg = i
    ori = {}
    i = 0
    while i < len(dataNodes):
        if dataNodes[i].input == 'ori':
            ori[i] = dataNodes[i].id
        i += 1
    return ori, schNodes
