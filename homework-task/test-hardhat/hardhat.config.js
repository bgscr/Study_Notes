require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config(); // 引入 dotenv
require("@openzeppelin/hardhat-upgrades");
require("@nomicfoundation/hardhat-verify");

const { ProxyAgent, setGlobalDispatcher } = require("undici");
const proxyAgent = new ProxyAgent("http://127.0.0.1:7078");
setGlobalDispatcher(proxyAgent);

module.exports = {
  solidity: "0.8.28",
  networks: {
    sepolia: {
      url: process.env.SEPOLIA_RPC_URL || "", // 从 .env 文件读取 URL
      accounts:
        process.env.PRIVATE_KEY !== undefined ? [process.env.PRIVATE_KEY] : [], // 从 .env 文件读取私钥
    },
  },
  etherscan: {
    apiKey: process.env.ETHERSCAN_API_KEY,
  },
};
