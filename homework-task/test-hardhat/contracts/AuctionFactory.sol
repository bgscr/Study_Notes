// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// ==========================================================================================
// 文件: contracts/AuctionFactory.sol
// 描述: 工厂合约，用于创建和管理所有拍卖合约的实例。采用 UUPS 可升级代理模式。
// ==========================================================================================
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "./Auction.sol";

contract AuctionFactory is Initializable, UUPSUpgradeable, OwnableUpgradeable {
    address public auctionImplementation;
    address[] public allAuctions;

    event AuctionCreated(
        address indexed auctionAddress,
        address indexed seller,
        address indexed nftContract,
        uint256 tokenId
    );

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(address _auctionImplementation) public initializer {
        __Ownable_init();
        auctionImplementation = _auctionImplementation;
    }

    function createAuction(
        address _nftContract,
        uint256 _tokenId,
        uint256 _startingBid,
        address _erc20Token, // address(0) for ETH
        address _ethUsdPriceFeed,
        address _erc20UsdPriceFeed
    ) external {
        // 1. 卖家将 NFT 转移给工厂
        IERC721(_nftContract).transferFrom(msg.sender, address(this), _tokenId);

        bytes memory data = abi.encodeWithSelector(
            Auction.initialize.selector,
            msg.sender,
            _nftContract,
            _tokenId,
            _startingBid,
            _erc20Token,
            _ethUsdPriceFeed,
            _erc20UsdPriceFeed
        );
        
        // 2. 创建一个新的代理合约来管理拍卖
        ERC1967Proxy proxy = new ERC1967Proxy(
            auctionImplementation,
            data
        );
        
        address auctionAddress = address(proxy);
        allAuctions.push(auctionAddress);

        // 3. 将 NFT 从工厂转移到新创建的拍卖合约
        IERC721(_nftContract).safeTransferFrom(address(this), auctionAddress, _tokenId);

        emit AuctionCreated(
            auctionAddress,
            msg.sender,
            _nftContract,
            _tokenId
        );
    }

    function setAuctionImplementation(address _newImplementation) external onlyOwner {
        auctionImplementation = _newImplementation;
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}