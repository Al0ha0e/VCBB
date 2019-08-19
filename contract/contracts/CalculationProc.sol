pragma solidity ^0.5.0;
contract CalculationProc {
    struct ansInfo{
        string ans;
        address payable participant;
    }
    ansInfo[] answerList;
    uint256 ansCnt;
    uint8 constant maxRewardedParticipantCount = 8;
    uint8 state;
    uint8 rewardedParticipantCount;
    uint256 constant minimumMasterFund=100 wei;
    uint256 constant minimumParticipantFund = 100 wei;
    uint256 startTime;
    uint256[maxRewardedParticipantCount] rewardDistribute;
    string jobID;
    address payable master;
    //mapping(address=>bool) visited;
    mapping(address=>bool) blackList;
    mapping(address=>uint256)funding;
    mapping(string=>address payable[]) result;
    event committed(address participant,string ansHash);
    event punished(address participant);
    event terminated(string ans,uint256 cnt);
    constructor(string memory id,
                uint256 st,
                uint256 fund,
                uint8 participantCount,
                uint256[maxRewardedParticipantCount] memory distribute)payable public{
        require(participantCount>0 && participantCount<=maxRewardedParticipantCount,"invalid rewardedParticipantCount");
        require(fund>=minimumMasterFund,"master fund not enough");
        uint256 tot = fund;
        for(uint8 i = 0;i<participantCount;i++) tot += distribute[i];
        require(tot==msg.value, "invalid fund");
        jobID = id;
        startTime = st;
        rewardedParticipantCount = participantCount;
        rewardDistribute = distribute;
        state = 1;
        master = msg.sender;
        funding[master] = fund;
    }
    function commit(string memory answerHash)public payable{
        require(msg.value>=minimumParticipantFund,"participant fund not enough");
        require(block.timestamp>=startTime,"preparing");
        require(state==1,"not running");
        if(blackList[msg.sender]) return;
        if(funding[msg.sender]!=0 || msg.sender==master){
            punish(msg.sender);
            return;
        }
        ansCnt++;
        answerList.push(ansInfo(answerHash,msg.sender));
        funding[msg.sender] = msg.value;
        emit committed(msg.sender,answerHash);
    }
    function punish(address participant) internal {
        if(blackList[participant]) return;
        blackList[participant] = true;
        emit punished(participant);
    }
    function terminate() public {
        require(block.timestamp>=startTime,"preparing");
        if(msg.sender!=master){
            punish(msg.sender);
            return;
        }
        uint256 maxcnt = 0;
        string memory maxans;
        for(uint256 i = 0;i<ansCnt;i++){
            address payable pt = answerList[i].participant;
            if(blackList[pt]) continue;
            result[answerList[i].ans].push(pt);
            uint256 l = result[answerList[i].ans].length;
            if(l>maxcnt){
                maxcnt = l;
                maxans = answerList[i].ans;
            }
        }
        uint256 l = result[maxans].length;
        if(l<rewardedParticipantCount){
            for(uint256 i = l;i<rewardedParticipantCount;i++){
                funding[master] += rewardDistribute[i];
            }
            rewardedParticipantCount = uint8(l);
        }
        uint256 i;
        for(i = 0;i<rewardedParticipantCount;i++){
            address payable pt = result[maxans][i];
            pt.send(rewardDistribute[i]+funding[pt]);
        }
        for(i = rewardedParticipantCount;i<l;i++){
            address payable pt = result[maxans][i];
            pt.send(funding[pt]);
        }
        if(!blackList[master]){
            master.send(funding[master]);
        }
        emit terminated(maxans,maxcnt);
    }
    function()external payable{}
}
