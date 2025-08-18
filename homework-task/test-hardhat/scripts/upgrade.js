// File: scripts/upgrade.js
const { ethers, upgrades } = require("hardhat");

// Your AuctionFactory proxy address
const PROXY_ADDRESS = "0x52B8Accad6219e05a0967B39cDfF9152277592D4";

async function main() {
    console.log("Preparing to upgrade the AuctionFactory contract...");

    // 1. Get the new V2 contract factory
    const AuctionFactoryV2 = await ethers.getContractFactory("AuctionFactoryV2");

    // 2. Call upgradeProxy
    // This will deploy the V2 implementation and update the proxy to point to it.
    const upgradedFactory = await upgrades.upgradeProxy(PROXY_ADDRESS, AuctionFactoryV2);
    await upgradedFactory.waitForDeployment();

    console.log("AuctionFactory has been successfully upgraded.");
    console.log("Proxy is at:", await upgradedFactory.getAddress());
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
