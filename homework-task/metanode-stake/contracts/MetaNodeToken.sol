// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// 这是一个标准的 ERC20 代币，增加了 Ownable 功能以便增发
contract MetaNodeToken is ERC20, Ownable {
    constructor(address initialOwner) ERC20("MetaNode Token", "MNT") Ownable(initialOwner) {
        // 在部署时，为您自己（初始所有者）铸造 100 万个代币
        _mint(initialOwner, 1000000 * 10**18);
    }

    // 允许所有者增发代币的函数
    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
    }
}