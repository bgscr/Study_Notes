// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.19;

import "forge-std/Test.sol";
// 导入两个版本的合约
import "../src/CalculatorV1.sol";
import "../src/CalculatorV2.sol";

contract CalculatorTest is Test {
    CalculatorV1 public calculatorV1;
    CalculatorV2 public calculatorV2;

    // setUp 函数会为每个测试用例运行一次
    function setUp() public {
        calculatorV1 = new CalculatorV1();
        calculatorV2 = new CalculatorV2();
    }

    // --- V1 Tests ---
    function test_V1_Add() public {
        calculatorV1.add(100, 50);
        assertEq(calculatorV1.lastResult(), 150, "V1 Add result should be 150");
    }

    function test_V1_Subtract() public {
        calculatorV1.subtract(100, 50);
        assertEq(calculatorV1.lastResult(), 50, "V1 Subtract result should be 50");
    }

    // --- V2 Tests ---
    function test_V2_Add() public view {
        uint256 result = calculatorV2.add(100, 50);
        assertEq(result, 150, "V2 Add result should be 150");
    }

    function test_V2_Subtract() public view {
        uint256 result = calculatorV2.subtract(100, 50);
        assertEq(result, 50, "V2 Subtract result should be 50");
    }
}
