async function main() {
  const [deployer] = await ethers.getSigners();

  console.log("正在用账户部署合约:", deployer.address);

  const MetaNodeToken = await ethers.getContractFactory("MetaNodeToken");
  // 部署时将合约所有权交给部署者
  const token = await MetaNodeToken.deploy(deployer.address);

  await token.waitForDeployment();

  console.log("MetaNodeToken 部署到了:", await token.getAddress());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});