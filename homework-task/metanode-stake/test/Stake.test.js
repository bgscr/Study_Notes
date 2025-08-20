const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("质押合约测试 (Stake Contract)", function () {
    let Stake, stake, owner, addr1, addr2;
    let rewardToken;

    // ... [部署模拟 ERC20 的函数和 beforeEach 钩子函数与之前相同，此处省略] ...
    async function deployERC20Mock(name, symbol, initialSupply) {
        const ERC20MockFactory = await ethers.getContractFactory("ERC20Mock");
        const token = await ERC20MockFactory.deploy(name, symbol, ethers.parseEther(initialSupply));
        await token.waitForDeployment();
        return token;
    }
    beforeEach(async function () {
        [owner, addr1, addr2] = await ethers.getSigners();
        rewardToken = await deployERC20Mock("MetaNode", "MN", "1000000");
        Stake = await ethers.getContractFactory("Stake");
        const metaNodePerBlock = ethers.parseEther("10");
        stake = await upgrades.deployProxy(Stake, [await rewardToken.getAddress(), metaNodePerBlock, owner.address], { initializer: 'initialize' });
        await stake.waitForDeployment();
        await rewardToken.transfer(await stake.getAddress(), ethers.parseEther("500000"));
    });


    describe("池管理 (Pool Management)", function () {
        it("应该成功添加一个用于原生货币的质押池", async function () {
            const nativeCurrencyAddress = "0x0000000000000000000000000000000000000000";
            await expect(stake.add(nativeCurrencyAddress, 100, ethers.parseEther("0.1"), 10))
                .to.not.be.reverted;
            const pool = await stake.poolInfo(0);
            expect(pool.stTokenAddress).to.equal(nativeCurrencyAddress);
            expect(pool.poolWeight).to.equal(100);
        });
    });

    describe("质押功能 (Staking)", function () {
        beforeEach(async function () {
            const nativeCurrencyAddress = "0x0000000000000000000000000000000000000000";
            await stake.add(nativeCurrencyAddress, 100, ethers.parseEther("0.1"), 10);
        });

        it("应该允许用户质押原生货币", async function () {
            const stakeAmount = ethers.parseEther("1");
            await expect(stake.connect(addr1).stake(0, stakeAmount, { value: stakeAmount }))
                .to.emit(stake, "Deposit")
                .withArgs(addr1.address, 0, stakeAmount);
            const userInfo = await stake.userInfo(0, addr1.address);
            expect(userInfo.stAmount).to.equal(stakeAmount);
            const poolInfo = await stake.poolInfo(0);
            expect(poolInfo.stTokenAmount).to.equal(stakeAmount);
        });

        it("如果质押数量小于最低要求，应该失败", async function () {
            const stakeAmount = ethers.parseEther("0.05");
            // [修正] 将期望的错误信息改为英文
            await expect(
                stake.connect(addr1).stake(0, stakeAmount, { value: stakeAmount })
            ).to.be.revertedWith("Stake: amount is less than minimum deposit");
        });

        it("如果原生货币质押时 msg.value 与 amount 不匹配，应该失败", async function () {
            const stakeAmount = ethers.parseEther("1");
            const incorrectValue = ethers.parseEther("0.5");
            // [修正] 将期望的错误信息改为英文
            await expect(
                stake.connect(addr1).stake(0, stakeAmount, { value: incorrectValue })
            ).to.be.revertedWith("Stake: msg.value must match _amount for native currency");
        });
    });
});