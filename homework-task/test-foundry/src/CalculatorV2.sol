// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.19;

/**
 * @title CalculatorV2
 * @dev 优化后的计算器合约。
 * 优化1: 函数通过 `returns` 直接返回值，而不是写入状态变量，避免了昂贵的 SSTORE 操作。
 * 优化2: 将函数可见性设为 `external` 并使用 `pure` 修饰符，因为它不读取也不修改状态。
 */
contract CalculatorV2 {
    function add(uint256 a, uint256 b) external pure returns (uint256) {
        return a + b;
    }

    function subtract(uint256 a, uint256 b) external pure returns (uint256) {
        return a - b;
    }
}
