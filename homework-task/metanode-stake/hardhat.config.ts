require("@nomicfoundation/hardhat-toolbox");
require("@openzeppelin/hardhat-upgrades");
require("dotenv").config();

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  // Update the solidity version here
  solidity: "0.8.24",
  networks: { // 添加 networks 配置
    sepolia: {
      url: process.env.SEPOLIA_RPC_URL || "",
      accounts: [process.env.PRIVATE_KEY, process.env.PRIVATE_KEY_USER].filter(key => key !== undefined),
    },
  },
  ignition: {
    enabled: false,
  },
};