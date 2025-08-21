// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.19;

/**
 * @title CalculatorV1
 * @dev 一个简单的计算器合约，包含加法和减法功能。
 * 包含一个状态变量 `lastResult` 用于存储上次计算的结果。
 * 注意：写入状态变量会消耗大量 Gas。
 */
contract CalculatorV1 {
    uint256 public lastResult;

    function add(uint256 a, uint256 b) public {
        uint256 result = a + b;
        lastResult = result; // 写入状态变量
    }

    function subtract(uint256 a, uint256 b) public {
        uint256 result = a - b;
        lastResult = result; // 写入状态变量
    }
}
