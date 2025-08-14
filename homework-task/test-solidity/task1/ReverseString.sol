// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ReversString{

    function reversString(string memory input) public pure returns (string memory output){
        bytes memory inBytes=bytes(input);
        for(uint i=0;i<inBytes.length/2;i++){
            bytes1 temp = inBytes[i];            
            inBytes[i]=inBytes[inBytes.length-i-1];
            inBytes[inBytes.length-i-1]=temp;
            
        }

        return string(inBytes);
    }
}