const { ethers, upgrades } = require("hardhat");

async function main() {
  // ########## 配置 ##########
  // 1. 在部署 MetaNodeToken 后，将得到的地址粘贴到这里
  const metaNodeTokenAddress = "0xCae846EAa76A67F623163775fd1D03113157f789"; // <---- 您的 MNT 代币地址

  // 2. 设置每个区块奖励的代币数量 (例如: 10个)
  const metaNodePerBlock = ethers.parseEther("10");
  // ##########################

  const [deployer] = await ethers.getSigners();
  console.log("正在用账户部署 Stake 合约:", deployer.address);

  if (!metaNodeTokenAddress || metaNodeTokenAddress.includes("...")) {
      console.error("错误：请先在脚本中设置 metaNodeTokenAddress！");
      return;
  }

  const Stake = await ethers.getContractFactory("Stake");

  console.log("正在部署 Stake 合约的代理...");
  
  const stake = await upgrades.deployProxy(
      Stake, 
      [metaNodeTokenAddress, metaNodePerBlock, deployer.address], 
      { initializer: 'initialize' }
  );

  await stake.waitForDeployment();
  const stakeProxyAddress = await stake.getAddress();

  console.log("Stake 合约 (代理) 部署到了:", stakeProxyAddress);
  
  // [修正] 使用新的 upgrades.erc1967.getImplementationAddress 函数
  const implementationAddress = await upgrades.erc1967.getImplementationAddress(
    stakeProxyAddress
  );
  console.log("Stake 实现合约部署到了:", implementationAddress);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});