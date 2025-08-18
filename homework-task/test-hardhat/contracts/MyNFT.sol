// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// 导入 OpenZeppelin 合约库
// ERC721URIStorage 包含了 ERC721 的所有功能，并增加了元数据存储能力
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
// Ownable 用于管理合约权限，确保只有合约所有者可以执行特定操作
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title MyNFT
 * @dev 一个更新后的基础 ERC721 NFT 合约，使用 uint256 来追踪 tokenID。
 */
contract MyNFT is ERC721URIStorage, Ownable {
    uint256 private _nextTokenId;

    constructor(
        address initialOwner
    ) ERC721("My NFT", "MNFT") Ownable(initialOwner) {}

    function mintNFT(
        address recipient,
        string memory tokenURI
    ) public onlyOwner returns (uint256) {
        uint256 newItemId = _nextTokenId;
        _nextTokenId++;

        _safeMint(recipient, newItemId);
        _setTokenURI(newItemId, tokenURI);

        return newItemId;
    }
}
