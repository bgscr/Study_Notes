// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title MergeSortedArrays
 * @dev 将两个已排序的 uint 数组合并为一个新的已排序数组。
 */
contract MergeSortedArrays {

    /**
     * @dev 将两个已排序的数组合并。
     * @param arr1 第一个已排序的数组。
     * @param arr2 第二个已排序的数组。
     * @return mergedArray 合并后的新排序数组。
     */
    function merge(uint[] memory arr1, uint[] memory arr2) public pure returns (uint[] memory) {
        uint i=0;
        uint j=0;
        uint k=0;

        uint len1=arr1.length;
        uint len2=arr2.length;

        uint[] memory mergedArray = new uint[](len1+len2);
        while(i<len1 && j<len2){
            if(arr1[i]<arr2[j]){
                mergedArray[k]=arr1[i];
                i++;
            }
            else{
                mergedArray[k]=arr2[j];
                j++;
            }
            k++;
        }

        while(i<len1){
                mergedArray[k]=arr1[i];
                i++;
                k++;
        }
        while(j<len2){
                mergedArray[k]=arr2[j];
                j++;
                k++;
        }

        return mergedArray;
    }
}