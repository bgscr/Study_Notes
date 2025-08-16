// test/Donation.test.js
const {
  loadFixture,
} = require("@nomicfoundation/hardhat-toolbox/network-helpers");
const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("DonationContract", function () {
  // 我们定义一个 "fixture" 函数来部署合约，这样可以复用部署逻辑
  async function deployDonationContractFixture() {
    // 获取签名者（钱包账户）
    const [owner, otherAccount] = await ethers.getSigners();

    // 部署合约，并将 owner 作为构造函数参数传入
    const DonationContractFactory = await ethers.getContractFactory("DonationContract");
    const donationContract = await DonationContractFactory.deploy(owner.address);
    
    // 返回需要用到的变量
    return { donationContract, owner, otherAccount };
  }

  // 测试用例 1: 检查部署是否成功
  it("Should set the right owner", async function () {
    const { donationContract, owner } = await loadFixture(deployDonationContractFixture);
    
    // 断言：合约的 owner() 函数返回的地址是否等于我们传入的 owner 地址
    expect(await donationContract.owner()).to.equal(owner.address);
  });

  // 测试用例 2: 检查捐赠功能
  it("Should allow users to donate and update contract balance", async function () {
    const { donationContract, otherAccount } = await loadFixture(deployDonationContractFixture);
    
    const donationAmount = ethers.parseEther("1.0"); // 捐赠 1 ETH

    // 让 otherAccount 发起一笔捐赠交易
    await donationContract.connect(otherAccount).donate({ value: donationAmount });

    // 断言：合约的余额是否等于捐赠的金额
    expect(await donationContract.getBalance()).to.equal(donationAmount);
  });
});