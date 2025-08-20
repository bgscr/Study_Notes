// contracts/test/ERC20Mock.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

// 这是一个用于测试的模拟 ERC20 合约
contract ERC20Mock is ERC20 {
    // 构造函数，在部署时初始化代币名称、符号和初始供应量
    constructor(string memory name, string memory symbol, uint256 initialSupply) ERC20(name, symbol) {
        // 将初始供应量铸造给部署者
        _mint(msg.sender, initialSupply);
    }
}