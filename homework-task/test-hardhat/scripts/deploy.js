const { ethers, upgrades } = require("hardhat");

async function main() {
    // --- 1. 获取部署者 ---
    const [deployer] = await ethers.getSigners();
    console.log("Deploying contracts with the account:", deployer.address);
    // **FIXED**: Updated to ethers v6+ syntax for getting balance
    const balance = await ethers.provider.getBalance(deployer.address);
    console.log("Account balance:", ethers.formatEther(balance));

    // --- 2. 部署 Auction 实现合约 ---
    const Auction = await ethers.getContractFactory("Auction");
    const auctionImplementation = await Auction.deploy();
    await auctionImplementation.waitForDeployment();
    const auctionImplementationAddress = await auctionImplementation.getAddress();
    console.log("Auction implementation deployed to:", auctionImplementationAddress);

    // --- 3. 部署 AuctionFactory 代理合约 ---
    // 使用 OpenZeppelin Upgrades 插件的 deployProxy 函数
    const AuctionFactory = await ethers.getContractFactory("AuctionFactory");
    const auctionFactory = await upgrades.deployProxy(AuctionFactory, [auctionImplementationAddress], {
        initializer: "initialize",
        kind: "uups", 
    });
    await auctionFactory.waitForDeployment();
    const auctionFactoryAddress = await auctionFactory.getAddress();
    console.log("AuctionFactory proxy deployed to:", auctionFactoryAddress);
    
    // --- 4. (可选) 部署 NFT 合约用于后续交互 ---
    const MyNFT = await ethers.getContractFactory("MyNFT");
    const myNFT = await MyNFT.deploy(deployer.address);
    await myNFT.waitForDeployment();
    const myNFTAddress = await myNFT.getAddress();
    console.log("MyNFT deployed to:", myNFTAddress);
    
    console.log("\n--- Deployment Summary ---");
    console.log(`Auction Logic: ${auctionImplementationAddress}`);
    console.log(`Auction Factory (Proxy): ${auctionFactoryAddress}`);
    console.log(`MyNFT: ${myNFTAddress}`);
    console.log("--------------------------\n");
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
