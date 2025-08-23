# **SPL 代币发行程序操作说明文档**

本文档详细介绍了如何部署和运行基于 Solana 和 Anchor 框架开发的 SPL 代币发行程序。该程序实现了代币的初始化、铸造和转移等核心功能。

### **一、 环境准备**

在开始之前，请确保您的开发环境中已安装以下工具。推荐使用 WSL2 (Windows) 或原生 Linux/macOS 环境。

1. **Rust & Cargo**: Solana 智能合约的开发语言和包管理器。  
   * 安装命令: curl \--proto '=https' \--tlsv1.2 \-sSf https://sh.rustup.rs | sh  
2. **Solana Tool Suite**: Solana 命令行工具集，用于与网络交互。  
   * 安装命令: curl \--proto '=https' \--tlsv1.2 \-sSfL https://solana-install.solana.workers.dev | bash  
3. **Anchor Framework**: Solana 智能合约开发框架。  
   * 安装命令: cargo install \--git https://github.com/coral-xyz/anchor avm \--locked \--force  
   * 安装后执行: avm install latest && avm use latest  
4. **Node.js & Yarn**: 用于运行测试脚本和客户端应用。  
   * 推荐使用 nvm 安装 Node.js: curl \-o- https://raw.githubusercontent.com/nvm-sh/nvm/master/install.sh | bash  
   * 安装 Node.js: nvm install \--lts  
   * 安装 Yarn: npm install \-g yarn

安装完成后，请重启终端并通过 solana \--version, anchor \--version, node \-v, yarn \--version 命令验证所有工具是否安装成功。

### **二、 项目编译与测试**

在部署之前，首先在本地验证程序是否能正常工作。

1. 克隆项目/进入项目目录  
   将本项目代码下载到本地，并在终端中进入项目根目录。  
2. 安装依赖  
   安装测试脚本所需的 JavaScript 依赖库，特别是 @solana/spl-token。  
   yarn add @solana/spl-token

3. 编译程序  
   该命令会将 programs/spl\_token 目录下的 Rust 代码编译成 BPF 字节码。  
   anchor build

   编译成功后，会在 target/ 目录下生成程序的 .so 文件和 IDL (.json) 文件。  
4. 运行本地测试  
   该命令会自动启动一个本地测试验证器，将程序部署上去，并执行 tests/ 目录下的测试脚本。  
   anchor test

   如果所有测试（初始化、铸造、转移）都显示为绿色通过，则证明程序逻辑正确。

### **三、 部署与运行**

#### **1\. 部署到本地网络 (Localnet)**

本地网络是在您自己电脑上运行的私有 Solana 集群，是开发和调试的首选。

* 启动本地验证器  
  打开一个新的终端窗口，运行以下命令来启动一个独立的本地测试网络：  
  solana-test-validator

  这个终端将持续输出日志，请保持其运行。  
* 部署程序  
  回到你原来的项目终端窗口，运行部署命令：  
  anchor deploy

  部署成功后，命令会输出你程序的 ID (Program ID)。请确保这个 ID 与 programs/spl\_token/src/lib.rs 文件中 declare\_id\! 的值一致。

#### **2\. 部署到开发网 (Devnet)**

开发网是一个公开的测试网络，可以模拟真实的网络环境。

* 配置网络  
  将 Solana CLI 的目标网络设置为开发网。  
  solana config set \--url devnet

* 获取测试 SOL  
  部署需要消耗 SOL 作为租金和交易费。你可以向开发网请求空投一些测试 SOL。  
  solana airdrop 2

  (可以请求 1-2 个 SOL，如果余额不足可以多执行几次)  
* 修改 Anchor.toml  
  打开项目根目录的 Anchor.toml 文件，将 \[provider\] 部分的 cluster 修改为 devnet。  
  \[provider\]  
  cluster \= "devnet"  
  wallet \= "/home/jayce/.config/solana/id.json"

* 部署到开发网  
  运行部署命令：  
  anchor deploy

#### **3\. 运行程序**

部署后，与程序的主要交互方式是通过客户端脚本（例如我们项目中的 tests/spl\_token.ts）。你可以基于这个测试脚本来开发自己的客户端应用或命令行工具，调用合约中的 initialize, mintTokens, transferTokens 等方法来发行和管理你的代币。