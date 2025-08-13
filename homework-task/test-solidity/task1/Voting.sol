// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title Voting
 * @dev 一个简单的投票合约，允许用户为候选人投票并查询票数。
 */
contract Voting {
    // 状态变量
    address public owner; // 合约的拥有者，只有他可以重置投票
    mapping(string => uint) public candidates; // 存储每个候选人得票数的映射
    string[] public candidateList; // 一个动态数组，用于追踪所有候选人的名字

    // 事件
    event Voted(address indexed voter, string indexed candidateName);
    event VotesReset();

    /**
     * @dev 构造函数，在合约部署时设置拥有者。
     */
    constructor() {
        owner = msg.sender;
    }

    /**
     * @dev 一个修饰器，限制只有合约拥有者才能执行某个函数。
     */
    modifier onlyOwner() {
        require(msg.sender == owner, "Only the owner can call this function.");
        _;
    }

    /**
     * @dev 允许用户为指定候选人投票。
     * @param _candidateName 要投票的候选人名字。
     */
    function vote(string memory _candidateName) public {
        // 检查候选人是否已在列表中，如果不在则添加
        if (!candidateExists(_candidateName)) {
            candidateList.push(_candidateName);
        }

        // 增加候选人的得票数
        candidates[_candidateName]++;

        // 触发投票事件
        emit Voted(msg.sender, _candidateName);
    }

    /**
     * @dev 获取某个候选人的当前得票数。
     * @param _candidateName 要查询的候选人名字。
     * @return 票数。
     */
    function getVotes(string memory _candidateName) public view returns (uint) {
        return candidates[_candidateName];
    }

    /**
     * @dev 重置所有候选人的得票数，只有合约拥有者可以调用。
     */
    function resetVotes() public onlyOwner {
        // 遍历所有候选人列表
        for (uint i = 0; i < candidateList.length; i++) {
            string memory candidateName = candidateList[i];
            // 将每个候选人的票数重置为0
            candidates[candidateName] = 0;
        }

        // 清空候选人列表数组
        delete candidateList;

        // 触发重置事件
        emit VotesReset();
    }

    /**
     * @dev 内部辅助函数，检查一个候选人是否存在于列表中。
     * @param _candidateName 要检查的候选人名字。
     * @return 如果存在则返回 true，否则返回 false。
     */
    function candidateExists(string memory _candidateName) internal view returns (bool) {
        for (uint i = 0; i < candidateList.length; i++) {
            // Solidity 中比较字符串需要对它们的哈希值进行比较
            if (keccak256(abi.encodePacked(candidateList[i])) == keccak256(abi.encodePacked(_candidateName))) {
                return true;
            }
        }
        return false;
    }
}