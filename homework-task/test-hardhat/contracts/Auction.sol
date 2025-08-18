// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/utils/ERC721HolderUpgradeable.sol";

contract Auction is Initializable, UUPSUpgradeable, ERC721HolderUpgradeable {
    address public seller;
    address public nftContract;
    uint256 public tokenId;
    uint256 public endTime;
    bool public ended;
    address public highestBidder;
    uint256 public highestBid;
    mapping(address => uint256) public bids;
    address public erc20Token;
    AggregatorV3Interface internal ethUsdPriceFeed;
    AggregatorV3Interface internal erc20UsdPriceFeed;

    event Bid(address indexed bidder, uint256 amount);
    event AuctionEnded(address winner, uint256 amount);
    event Withdrawal(address bidder, uint256 amount);

    modifier notEnded() {
        require(block.timestamp < endTime, "Auction already ended");
        _;
    }

    modifier onlySeller() {
        require(msg.sender == seller, "Only seller can call this");
        _;
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
        address _seller,
        address _nftContract,
        uint256 _tokenId,
        uint256 _startingBid,
        address _erc20Token,
        address _ethUsdPriceFeed,
        address _erc20UsdPriceFeed
    ) public initializer {
        __UUPSUpgradeable_init();
        __ERC721Holder_init();
        seller = _seller;
        nftContract = _nftContract;
        tokenId = _tokenId;
        endTime = block.timestamp + 7 days;
        ended = false;
        highestBid = _startingBid;
        erc20Token = _erc20Token;
        if (_ethUsdPriceFeed != address(0)) {
            ethUsdPriceFeed = AggregatorV3Interface(_ethUsdPriceFeed);
        }
        if (_erc20UsdPriceFeed != address(0)) {
            erc20UsdPriceFeed = AggregatorV3Interface(_erc20UsdPriceFeed);
        }
    }

    function bid() external payable notEnded {
        if (!isEthAuction()) {
            revert("ERC20 bids need a separate function with amount parameter");
        }
        require(msg.value > highestBid, "Bid must be higher than current highest bid");
        if (highestBidder != address(0)) {
            (bool success, ) = payable(highestBidder).call{value: highestBid}("");
            require(success, "Failed to refund previous bidder");
        }
        highestBidder = msg.sender;
        highestBid = msg.value;
        bids[msg.sender] = msg.value;
        emit Bid(msg.sender, msg.value);
    }

    function bidWithErc20(uint256 _amount) external notEnded {
        require(!isEthAuction(), "This auction is for ETH");
        require(_amount > highestBid, "Bid must be higher than current highest bid");
        if (highestBidder != address(0)) {
            IERC20(erc20Token).transfer(highestBidder, highestBid);
        }
        IERC20(erc20Token).transferFrom(msg.sender, address(this), _amount);
        highestBidder = msg.sender;
        highestBid = _amount;
        bids[msg.sender] = _amount;
        emit Bid(msg.sender, _amount);
    }

    function endAuction() external {
        require(block.timestamp >= endTime, "Auction not yet ended");
        require(!ended, "Auction already ended");
        ended = true;
        if (highestBidder != address(0)) {
            IERC721(nftContract).safeTransferFrom(address(this), highestBidder, tokenId);
            if (isEthAuction()) {
                (bool success, ) = payable(seller).call{value: highestBid}("");
                require(success, "Failed to send funds to seller");
            } else {
                IERC20(erc20Token).transfer(seller, highestBid);
            }
            emit AuctionEnded(highestBidder, highestBid);
        } else {
            IERC721(nftContract).safeTransferFrom(address(this), seller, tokenId);
            emit AuctionEnded(address(0), 0);
        }
    }

    function isEthAuction() public view returns (bool) {
        return erc20Token == address(0);
    }

    function getLatestPrice() public view returns (int) {
        if (isEthAuction()) {
            require(address(ethUsdPriceFeed) != address(0), "ETH price feed not set");
            (, int price, , , ) = ethUsdPriceFeed.latestRoundData();
            return price;
        } else {
            require(address(erc20UsdPriceFeed) != address(0), "ERC20 price feed not set");
            (, int price, , , ) = erc20UsdPriceFeed.latestRoundData();
            return price;
        }
    }

    function getBidInUSD(uint256 _amount) public view returns (uint256) {
        int price = getLatestPrice();
        return (_amount * uint256(price)) / (10**26);
    }

    function _authorizeUpgrade(address /* newImplementation */) internal view override {
        require(msg.sender == seller, "Only seller can upgrade");
    }
}