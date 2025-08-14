// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BinarySearch {
    function binarySearch(uint256[] memory array, uint256 target)
        public
        pure
        returns (uint256)
    {
        uint256 i = 0;
        uint256 j = array.length - 1;

        while (i <= j) {
            uint256 mid = i + (j - i) / 2;
            if (target == array[mid]) {
                return mid;
            }
            if (target < array[mid]) {
                if (mid == 0) {
                    revert("no target num");
                }
                j = mid - 1;
            }
            if (target > array[mid]) {
                i = mid + 1;
            }
        }

        revert("no target num");
    }
}
