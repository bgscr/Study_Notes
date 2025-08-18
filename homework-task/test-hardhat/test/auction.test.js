const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("NFT Auction Marketplace", function () {
    let auctionFactory, auction, myNFT, mockERC20;
    let owner, seller, bidder1, bidder2;
    let ethUsdPriceFeed, erc20UsdPriceFeed; // Declared here to be accessible everywhere in the describe block

    beforeEach(async function () {
        [owner, seller, bidder1, bidder2] = await ethers.getSigners();

        // 部署 Mock 预言机 and assign to the higher-scoped variables
        const MockV3Aggregator = await ethers.getContractFactory("MockV3Aggregator");
        ethUsdPriceFeed = await MockV3Aggregator.deploy(8, 2000 * 10**8);
        await ethUsdPriceFeed.waitForDeployment();
        erc20UsdPriceFeed = await MockV3Aggregator.deploy(8, 1 * 10**8);
        await erc20UsdPriceFeed.waitForDeployment();

        // 部署 NFT 合约
        const MyNFT = await ethers.getContractFactory("MyNFT");
        myNFT = await MyNFT.deploy(owner.address);
        await myNFT.waitForDeployment();

        // 部署 Mock ERC20 合约
        const MockERC20Factory = await ethers.getContractFactory("MockERC20");
        mockERC20 = await MockERC20Factory.deploy();
        await mockERC20.waitForDeployment();
        
        await mockERC20.transfer(bidder1.address, ethers.parseEther("1000"));
        await mockERC20.transfer(bidder2.address, ethers.parseEther("1000"));

        // 部署 Auction 实现合约
        const Auction = await ethers.getContractFactory("Auction");
        const auctionImplementation = await Auction.deploy();
        await auctionImplementation.waitForDeployment();

        // 部署 AuctionFactory (代理)
        const AuctionFactory = await ethers.getContractFactory("AuctionFactory");
        auctionFactory = await upgrades.deployProxy(AuctionFactory, [await auctionImplementation.getAddress()], {
            initializer: "initialize",
            kind: "uups",
        });
        await auctionFactory.waitForDeployment();
    });

    // This helper function now correctly uses the variables from the beforeEach hook
    async function createAuction(isEth = true) {
        await myNFT.connect(owner).mintNFT(seller.address, "");
        const tokenId = 0;

        await myNFT.connect(seller).approve(await auctionFactory.getAddress(), tokenId);

        const startingBid = ethers.parseEther("0.1");
        const erc20TokenAddress = isEth ? ethers.ZeroAddress : await mockERC20.getAddress();

        const tx = await auctionFactory.connect(seller).createAuction(
            await myNFT.getAddress(),
            tokenId,
            startingBid,
            erc20TokenAddress,
            await ethUsdPriceFeed.getAddress(), // Correctly access the address
            await erc20UsdPriceFeed.getAddress() // Correctly access the address
        );
        
        const receipt = await tx.wait();
        // **FIXED**: Correctly parse event logs for ethers v6+
        const event = receipt.logs.find(log => log.fragment && log.fragment.name === 'AuctionCreated');
        const auctionAddress = event.args.auctionAddress;
        
        auction = await ethers.getContractAt("Auction", auctionAddress);
        return { tokenId, startingBid };
    }

    it("Should create an ETH auction successfully", async function () {
        await createAuction(true);
        expect(await auction.seller()).to.equal(seller.address);
        expect(await auction.isEthAuction()).to.be.true;
        expect(await myNFT.ownerOf(0)).to.equal(await auction.getAddress());
    });

    it("Should handle the full ETH auction flow", async function () {
        const { tokenId } = await createAuction(true);

        const bid1Amount = ethers.parseEther("0.2");
        await auction.connect(bidder1).bid({ value: bid1Amount });
        
        const bid2Amount = ethers.parseEther("0.3");
        await expect(auction.connect(bidder2).bid({ value: bid2Amount }))
            .to.changeEtherBalance(bidder1, bid1Amount);

        await ethers.provider.send("evm_increaseTime", [7 * 24 * 60 * 60 + 1]);
        await ethers.provider.send("evm_mine");

        await expect(auction.connect(seller).endAuction())
            .to.changeEtherBalance(seller, bid2Amount);

        expect(await myNFT.ownerOf(tokenId)).to.equal(bidder2.address);
    });

    it("Should handle the full ERC20 auction flow", async function () {
        const { tokenId } = await createAuction(false);

        const bid1Amount = ethers.parseEther("100");
        await mockERC20.connect(bidder1).approve(await auction.getAddress(), bid1Amount);
        await auction.connect(bidder1).bidWithErc20(bid1Amount);
        
        const bid2Amount = ethers.parseEther("120");
        await mockERC20.connect(bidder2).approve(await auction.getAddress(), bid2Amount);
        
        await expect(auction.connect(bidder2).bidWithErc20(bid2Amount))
            .to.changeTokenBalance(mockERC20, bidder1, bid1Amount);

        await ethers.provider.send("evm_increaseTime", [7 * 24 * 60 * 60 + 1]);
        await ethers.provider.send("evm_mine");

        await expect(auction.connect(seller).endAuction())
            .to.changeTokenBalance(mockERC20, seller, bid2Amount);
        
        expect(await myNFT.ownerOf(tokenId)).to.equal(bidder2.address);
    });
    
    it("Should allow the owner to upgrade the factory contract", async function () {
        const AuctionFactoryV2 = await ethers.getContractFactory("AuctionFactory");
        const implementationV2 = await AuctionFactoryV2.deploy();
        await implementationV2.waitForDeployment();

        await auctionFactory.connect(owner).upgradeTo(await implementationV2.getAddress());
        
        const implementationAddress = await upgrades.erc1967.getImplementationAddress(await auctionFactory.getAddress());
        expect(implementationAddress).to.equal(await implementationV2.getAddress());
    });
});