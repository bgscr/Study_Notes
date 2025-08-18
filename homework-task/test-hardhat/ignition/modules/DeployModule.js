// 文件: ignition/modules/DeployModule.js
const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

module.exports = buildModule("AuctionModule", (m) => {
  // --- 1. 部署 Auction 实现合约 ---
  const auctionImplementation = m.contract("Auction");

  // --- 2. 部署 AuctionFactory 代理合约 ---
  // 使用 m.contract.proxy() 来部署可升级的代理合约
  const auctionFactory = m.contract.proxy("AuctionFactory", [auctionImplementation], {
    initializer: "initialize",
  });

  // --- 3. (可选) 部署 NFT 合约用于测试 ---
  // 注意：部署测试用的 NFT 和 ERC20 通常在测试脚本中完成，
  // 实际部署时可能不需要这一步，除非你的 DApp 需要一个初始的 NFT 集合。
  const owner = m.getAccount(0); // 获取部署者账户
  const myNFT = m.contract("MyNFT", [owner]);

  console.log("Deployment setup complete.");
  console.log("Auction Logic Implementation will be deployed.");
  console.log("Auction Factory (Proxy) will be deployed.");
  console.log("MyNFT (for testing) will be deployed.");

  return { auctionImplementation, auctionFactory, myNFT };
});
