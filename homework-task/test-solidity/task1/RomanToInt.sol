// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract RomanToInt{

    function romanToInt(string memory s) public pure returns(int result){


        bytes memory input = bytes(s);
        uint length = input.length;
        int r=0;
        for (uint i=0;i<length-1;i++){
            int current=hepler(input[i]);
            int next=hepler(input[i+1]);

            if(current<next){
                r-=current;
            }else{
                r+=current;
            }
        }

        r+=hepler(input[length-1]);
        return r;
    }


    function hepler(bytes1 v) internal pure returns (int result){
        if (v=='I') return 1;
        if (v=='V') return 5;
        if (v=='X') return 10;
        if (v=='L') return 50;
        if (v=='C') return 100 ;
        if (v=='D') return 500;
        if (v=='M') return 1000;

        revert("no support character");
    }
}