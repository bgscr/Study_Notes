// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

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

    /**
     * @dev 合约的构造函数.
     * @param initialOwner 合约的初始所有者地址.
     */
    constructor(address initialOwner)
        ERC721("My NFT", "MNFT") // 设置 NFT 的名称和符号
        Ownable(initialOwner) // 将部署合约的地址设为所有者
    {}

    /**
     * @dev 铸造一个新的 NFT.
     * @param recipient 将接收新 NFT 的地址.
     * @param tokenURI 新 NFT 的元数据链接.
     * @return 返回新铸造的 NFT 的 ID.
     */
    function mintNFT(address recipient, string memory tokenURI)
        public
        onlyOwner // 限制只有合约所有者才能调用此函数
        returns (uint256)
    {
        // 获取当前的 token ID 用于铸造，然后将计数器加一
        uint256 newItemId = _nextTokenId;
        _nextTokenId++;

        _safeMint(recipient, newItemId);
        _setTokenURI(newItemId, tokenURI);

        return newItemId;
    }
}