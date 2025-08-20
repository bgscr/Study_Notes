const { ethers } = require("hardhat");

async function main() {
  // ====================== 配置区 ======================
  // ！！！注意：请在这里粘贴您刚刚重新部署后得到的【新】Stake 合约地址 ！！！
  const stakeContractAddress = "0x4f9b7b856126eEAEd56924116656754dDf9b07db"; // <---- 替换为您的【新】Stake 合约地址
  const mntContractAddress = "0xCae846EAa76A67F623163775fd1D03113157f789"; // MNT 地址不变

  const [owner, user] = await ethers.getSigners();
  console.log(`合约拥有者 (Owner): ${owner.address}`);
  console.log(`普通用户 (User):  ${user.address}`);
  console.log("----------------------------------------------------");
  // ====================================================

  const stakeContract = await ethers.getContractAt("Stake", stakeContractAddress);
  const mntContract = await ethers.getContractAt("MetaNodeToken", mntContractAddress);
  console.log("✅ 成功获取 Stake 和 MetaNodeToken 合约实例");

  // --- 省略了前面的步骤，因为您已经成功执行过了 ---
  // ... 如果需要可以取消注释来重新执行 ...
  
  // --- 检查奖励 ---
  console.log("\n💰 正在调用 pendingReward 函数检查奖励...");
  const pendingReward = await stakeContract.pendingReward(0, user.address);
  console.log(`   USER 在池 0 的待领取奖励: ${ethers.formatEther(pendingReward)} MNT`);
  console.log("----------------------------------------------------");

  // --- 领取奖励 (已修正) ---
  console.log("\nUSER: 正在尝试领取奖励...");
  if (pendingReward > 0) {
    try {
      const tx = await stakeContract.connect(user).claimReward(0);
      
      // [改进1] 立即打印交易哈希
      console.log(`   ... 交易已发送! 哈希: ${tx.hash}`);
      console.log("   ... 正在等待交易确认, 这可能需要一些时间...");
      
      // 等待交易被确认
      await tx.wait();

      const userMntBalance = await mntContract.balanceOf(user.address);
      console.log(`\n✅ USER 领取奖励成功!`);
      console.log(`   USER 现在的 MNT 余额: ${ethers.formatEther(userMntBalance)}`);

    } catch (error) {
      // [改进2] 捕获并打印详细错误
      console.error("\n❌ 领取奖励失败! 错误原因:");
      console.error(error);
    }
  } else {
    console.log("ℹ️  没有可领取的奖励。");
  }
  console.log("----------------------------------------------------");
  console.log("\n🎉 最终版交互脚本执行完毕!");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});