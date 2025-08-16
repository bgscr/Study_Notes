// contracts/DonationContract.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/access/Ownable.sol";

contract DonationContract is Ownable {
    constructor(address initialOwner) Ownable(initialOwner) {}

    function donate() public payable {
        require(msg.value > 0, "Donation must be > 0");
    }

    function withdraw(address payable _to) public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        
        (bool success, ) = _to.call{value: balance}("");
        require(success, "Withdraw failed");
    }

    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
}