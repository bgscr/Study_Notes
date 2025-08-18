// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "./Auction.sol";

// We are NOT inheriting PausableUpgradeable anymore to maintain storage compatibility.
contract AuctionFactoryV2 is Initializable, UUPSUpgradeable, OwnableUpgradeable {
    // --- V1 State Variables (Order MUST NOT change) ---
    address public auctionImplementation;
    address[] public allAuctions;

    // --- V2 State Variables (New variables MUST be added at the end) ---
    bool private _paused;

    event AuctionCreated(
        address indexed auctionAddress,
        address indexed seller,
        address indexed nftContract,
        uint256 tokenId
    );
    
    // --- V2 Events ---
    event Paused(address account);
    event Unpaused(address account);


    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(address _auctionImplementation) public initializer {
        __Ownable_init();
        __UUPSUpgradeable_init();
        // We do not initialize Pausable here anymore.
        auctionImplementation = _auctionImplementation;
    }

    // --- V2 Functions (Manually implemented from Pausable) ---
    modifier whenNotPaused() {
        require(!paused(), "Pausable: paused");
        _;
    }

    function paused() public view returns (bool) {
        return _paused;
    }

    function _pause() internal whenNotPaused {
        _paused = true;
        emit Paused(owner());
    }

    function _unpause() internal {
        require(paused(), "Pausable: not paused");
        _paused = false;
        emit Unpaused(owner());
    }

    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }
    // --- End of V2 Functions ---


    function createAuction(
        address _nftContract,
        uint256 _tokenId,
        uint256 _startingBid,
        address _erc20Token,
        address _ethUsdPriceFeed,
        address _erc20UsdPriceFeed
    ) external whenNotPaused { // Added the whenNotPaused modifier
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

        ERC1967Proxy proxy = new ERC1967Proxy(
            auctionImplementation,
            data
        );
        
        address auctionAddress = address(proxy);
        allAuctions.push(auctionAddress);

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