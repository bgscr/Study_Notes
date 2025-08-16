// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title BeggingContract
 * @dev 增加了事件、排行榜和时间限制功能的最终版捐赠合约。
 */
contract BeggingContract is Ownable {

    // --- 核心状态变量 ---
    mapping(address => uint256) public donations; // 记录每个地址的捐赠总额
    
    // --- 排行榜功能所需变量 ---
    address[] private donors; // 存储所有不重复的捐贈者地址
    mapping(address => bool) private hasDonated; // 用于快速检查某地址是否已捐赠过

    // --- 时间限制功能所需变量 ---
    uint256 public immutable startTime; // 捐赠开始时间
    uint256 public immutable endTime;   // 捐赠结束时间

    // --- 事件定义 ---
    event Donation(address indexed donor, uint256 amount);

    /**
     * @dev 合约的构造函数.
     * @param _initialOwner 合约的初始所有者地址.
     * @param _durationInSeconds 捐赠活动允许的持续时间 (以秒为单位).
     */
    constructor(address _initialOwner, uint256 _durationInSeconds) Ownable(_initialOwner) {
        require(_durationInSeconds > 0, "Duration must be greater than zero");
        startTime = block.timestamp;
        endTime = block.timestamp + _durationInSeconds;
    }

    /**
     * @dev 接受捐赠的核心函数.
     */
    function donate() public payable {
        // 挑战 #3: 检查是否在捐赠时间窗口内
        require(block.timestamp >= startTime && block.timestamp <= endTime, "Donation period is over or not started yet");
        
        require(msg.value > 0, "Donation must be greater than zero");

        // 更新捐赠总额
        donations[msg.sender] += msg.value;

        // 挑战 #2: 如果是新捐赠者，将其添加到捐赠者列表中
        if (!hasDonated[msg.sender]) {
            hasDonated[msg.sender] = true;
            donors.push(msg.sender);
        }

        // 挑战 #1: 触发捐赠事件
        emit Donation(msg.sender, msg.value);
    }

    /**
     * @dev 合约所有者提取合约中的所有资金.
     */
    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        
        (bool success, ) = owner().call{value: balance}("");
        require(success, "Failed to send Ether");
    }

    /**
     * @dev 查询指定地址的总捐赠金额.
     */
    function getDonation(address _donor) public view returns (uint256) {
        return donations[_donor];
    }
    
    /**
     * @dev 查询合约当前的总余额.
     */
    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
    
    /**
     * @dev 挑战 #2: 获取捐赠金额最高的前3个地址.
     * @return 返回一个包含最多3个地址的数组.
     */
    function getTopDonors() public view returns (address[3] memory) {
        address[3] memory topAddresses;
        uint256[3] memory topAmounts;

        for (uint i = 0; i < donors.length; i++) {
            address currentDonor = donors[i];
            uint256 currentDonation = donations[currentDonor];

            // 比较并插入到排行榜
            if (currentDonation > topAmounts[0]) {
                topAmounts[2] = topAmounts[1];
                topAddresses[2] = topAddresses[1];
                topAmounts[1] = topAmounts[0];
                topAddresses[1] = topAddresses[0];
                topAmounts[0] = currentDonation;
                topAddresses[0] = currentDonor;
            } else if (currentDonation > topAmounts[1]) {
                topAmounts[2] = topAmounts[1];
                topAddresses[2] = topAddresses[1];
                topAmounts[1] = currentDonation;
                topAddresses[1] = currentDonor;
            } else if (currentDonation > topAmounts[2]) {
                topAmounts[2] = currentDonation;
                topAddresses[2] = currentDonor;
            }
        }
        return topAddresses;
    }
}