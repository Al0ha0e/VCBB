pragma solidity ^0.5.0;

contract CalculationProc {
    address payable foundation = address(0x583031D1113aD414F02576BD6afaBfb302140225);
    address payable master;//creater of the contract
    address payable[] participantList;//store all of the participants
    uint8 contractState;// 1 started 2 terminated 
    uint8 constant maxParticipantCount = 8;
    uint8 participantCount;//max number of participants restricted by master
    uint256 startTime;
    uint256 constant minimumParticipantFund = 100 wei;
    uint256 constant minimumMasterFund = 100 wei;
    uint256[2] feeFrac = [1,100];
    uint256[maxParticipantCount+2] fundDistribute;
    mapping(address => uint256)  participantFund;
    mapping(string => address payable[maxParticipantCount]) participants;// mapping ans to participants
    mapping(address => string) participantAnswer;
    mapping(string => uint8) answerCnt;
    mapping(address => bool) blackList;
    event CalcContractCreated(address _master,uint256 _fund,uint8 _participantCount);
    event CalcContractTerminated(uint256 _time);
    event Punish(address _lier);

    constructor(uint256 st,uint8 numParticipants,uint256[maxParticipantCount+2] memory distribute)payable public{
        require(msg.value*feeFrac[0]/feeFrac[1]<=distribute[maxParticipantCount],"fee not enough");
        require(distribute[maxParticipantCount+1]>=minimumMasterFund,"fund not enough");
        require(numParticipants>0,"no participant");
        require(numParticipants<=maxParticipantCount,"too much participants");
        uint256 sum = 0;
        for(uint256 i = 0;i<=maxParticipantCount+1;i++)sum += distribute[i];
        require(sum==msg.value,"fund not equals to expected");
        uint256 t = fundDistribute[maxParticipantCount];
        fundDistribute[maxParticipantCount] = 0;
        foundation.transfer(t);
        master = msg.sender;
        participantCount = numParticipants;
        startTime = st;
        fundDistribute = distribute;
        participantFund[master] = distribute[maxParticipantCount+1];
        contractState = 1;//started
        emit CalcContractCreated(msg.sender,msg.value,numParticipants);
    }
    function terminate(string memory answer) public {
        require(contractState==1,"state not started");
        if(msg.sender!=master) {
            punish(msg.sender);
            return;
        }
        contractState = 2;
        uint8 cnt = answerCnt[answer];
        uint256 i;
        uint256 retnum = 0;
        address payable pt;
        for(i = 0;i<cnt;i++){
            pt = participants[answer][i];
            retnum = fundDistribute[i]+participantFund[pt];
            fundDistribute[i] = 0;
            participantFund[pt] = 0;
            pt.send(retnum);
        }
        retnum = 0;
        for(;i<participantCount;i++) retnum += fundDistribute[i];
        retnum += participantFund[master];
        master.send(retnum);
        emit CalcContractTerminated(block.timestamp);
        selfdestruct(foundation);
    }
    function commit(string memory answer)public payable{
        require(block.timestamp>=startTime,"invalid time");
        require(msg.value>=minimumParticipantFund,"fund less than the minimum");
        require(contractState==1,"contract not open");
        require(!blackList[msg.sender],"participant in black list");
        uint8 cnt = answerCnt[answer];
        require(cnt<participantCount,"too much participants");
        if(msg.sender==master) {
            punish(master);
            return;
        }
        string memory s = participantAnswer[msg.sender];
        if(bytes(s).length>0){
            punish(msg.sender);
            for(uint256 i = 0;i<cnt;i++){
                if(msg.sender==participants[s][i]) {
                    for(i++;i<cnt;i++) participants[s][i-1] = participants[s][i];
                    delete(participants[s][cnt-1]);
                    answerCnt[s]--;
                    return;
                }
            }
        }
        participants[answer][cnt] = msg.sender;
        answerCnt[answer]++;
        participantList.push(msg.sender);
        participantFund[msg.sender] = msg.value;
    }
    function punish(address participant) internal {
        if(blackList[participant]) return;
        emit Punish(participant);
        blackList[participant] = true;
        uint256 t = participantFund[participant];
        if(t==0) return;
        participantFund[participant] = 0;
        foundation.transfer(t);// foundation.punish()
    }
    function()external payable{}
}
