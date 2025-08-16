require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config(); // 引入 dotenv

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.28",
  networks: {
    sepolia: {
      url: process.env.SEPOLIA_RPC_URL || "", // 从 .env 文件读取 URL
      accounts:
        process.env.PRIVATE_KEY !== undefined ? [process.env.PRIVATE_KEY] : [], // 从 .env 文件读取私钥
    },
  },
};
