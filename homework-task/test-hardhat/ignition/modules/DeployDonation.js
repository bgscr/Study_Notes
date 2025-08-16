// ignition/modules/DeployDonation.js
const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

// 这是一个 Ignition 模块，定义了部署合约的计划
module.exports = buildModule("DonationModule", (m) => {
  // 从部署者账户中获取地址作为合约的初始所有者
  const owner = m.getAccount(0);

  // 'm.contract()' 是核心函数，告诉 Ignition 我们要部署一个合约
  // 第一个参数 "DonationContract" 必须与你的合约名完全一致
  const donationContract = m.contract("DonationContract", {
    // 第二个参数是构造函数所需的参数数组
    // 我们的 DonationContract 的 constructor 需要一个 initialOwner 地址
    args: [owner], 
  });

  // 模块返回一个对象，键是我们部署的合约，值是合约实例
  return { donationContract };
});