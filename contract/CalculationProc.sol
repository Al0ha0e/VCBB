pragma solidity ^0.5.0;

contract CalculationProc {
    struct ansInfo{
        uint32 ans;
        address payable participant;
    }
    ansInfo[] answerList;
    uint8 ansCnt;
    uint8 constant maxRewardedParticipantCount = 8;
    uint8 state;
    uint8 rewardedParticipantCount;
    uint256 constant minimumMasterFund=1000000 wei;
    uint256 constant minimumParticipantFund = 1000 wei;
    uint256 startTime;
    uint256[maxRewardedParticipantCount] rewardDistribute;
    string jobID;
    address payable master;
    mapping(address=>bool) visited;
    mapping(address=>bool) blackList;
    constructor(string memory id,uint256 st,uint8 participantCount,uint256[maxRewardedParticipantCount] memory distribute)payable public{
        require(msg.value>minimumMasterFund,"master fund not enough");
        require(participantCount>0 && participantCount<=maxRewardedParticipantCount,"invalid rewardedParticipantCount");
        jobID = id;
        startTime = st;
        rewardedParticipantCount = participantCount;
        rewardDistribute = distribute;
        state = 1;
        master = msg.sender;
    }
    function terminate(string memory answer) public {
    }
    function commit(uint32 answer)public payable{
        require(msg.value>=minimumParticipantFund,"participant fund not enough");
        require(block.timestamp>=startTime,"preparing");
        require(state==1,"not running");
        if(blackList[msg.sender]) return;
        if(visited[msg.sender] || msg.sender==master){
            punish(msg.sender);
            return;
        }
        ansCnt++;
        answerList[ansCnt].ans = answer;
        answerList[ansCnt].participant = msg.sender;
        visited[msg.sender] = true;
    }
    function punish(address participant) internal {
    }
    function()external payable{}
}
