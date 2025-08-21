# Foundry 学习笔记：核心组件与 Windows 环境配置要点

通过这次 Gas 优化任务，把 Foundry 的安装、配置和使用流程完整地走了一遍。这份笔记主要记录了 Foundry 框架的核心概念，以及在 Windows 上配置和使用时的一些关键点和踩过的坑。

## 一、Foundry 框架核心组件

Foundry 是一个用 Rust 写的以太坊开发工具，速度很快。它主要由四个工具组成：

### 1. Forge: 测试框架
* **核心优势**: 最大的好处是能**直接用 Solidity 写测试**，不用像 Hardhat 那样来回切换 JS/TS，思路更连贯。
* **关键功能**: 自带的 Gas 报告 (`forge test --gas-report`) 非常好用，做合约性能分析和优化基本就靠它。
* **常用命令**:
    * `forge build`: 编译合约。
    * `forge test`: 跑测试，可以加上 `-vvv` 看更详细的日志。
    * `forge test --gas-report`: 跑测试并输出 Gas 报告。
    * `forge coverage`: 生成代码覆盖率报告。
    * `forge create <合约名> --private-key <私钥>`: 部署合约。

### 2. Cast: 命令行工具 (CLI)
* **核心优势**: 一个强大的命令行工具，能直接和以太坊网络交互，查数据、发交易都很方便，适合快速验证一些小想法。
* **常用命令**:
    * `cast call <合约地址> "functionName(uint256)" 123 --rpc-url <RPC地址>`: 调用合约的 view/pure 函数。
    * `cast send <合约地址> "functionName(uint256)" 123 --private-key <私钥> --rpc-url <RPC地址>`: 发送交易来调用写函数。
    * `cast balance <钱包地址>`: 查钱包余额。
    * `cast block latest`: 获取最新区块信息。
    * `cast --to-wei 1 ether`: 单位换算，ether 转 wei。

### 3. Anvil: 本地测试节点
* **核心优势**: 启动飞快，一个本地的以太坊节点，自带测试账户和 ETH，用于开发和测试。
* **常用命令**:
    * `anvil`: 启动默认的本地节点。
    * `anvil --fork-url <主网或测试网的RPC地址>`: 启动一个分叉节点，可以直接在本地模拟主网环境进行测试，这个功能非常实用。
    * `anvil -b 5`: 设置每 5 秒自动挖一个块。

### 4. Chisel: Solidity REPL (交互式解释器)
* **核心优势**: 一个 Solidity 的命令行交互环境，想测试一小段代码逻辑时，不用专门去写个合约和测试文件，直接在里面敲就行。
* **常用命令**:
    * `chisel`: 启动交互环境。
    * 在 Chisel 环境内:
        ```solidity
        > uint256 a = 100;
        > uint256 b = 50;
        > a + b
        150
        > bytes.concat("hello", " world")
        "hello world"
        ```

## 二、Windows 环境使用 Foundry 的要点

在 Windows 上用 Foundry，关键就是**所有操作都在 WSL (Windows Subsystem for Linux) 里进行**。

### 1. 安装与环境：必须在 WSL 内部

* **安装方式**: 只能在 WSL (推荐 Ubuntu) 里装 Foundry。直接在 Windows PowerShell 或 CMD 里是装不上的。
* **启动环境**: 最好是从 Windows 开始菜单直接搜 "Ubuntu" 应用来启动终端。这样能保证是以自己的普通用户身份登录。如果是在 CMD 里直接敲 `wsl`，有时候会默认用 `root` 用户登录，导致找不到之前装好的 `foundryup` 命令。

### 2. 开发工作流：VS Code + WSL 插件

* **代码编辑**: 在 Windows 上装 VS Code，然后在 VS Code 的扩展商店里搜并安装官方的 **WSL** 插件。
* **打开项目**: 在 Ubuntu 终端里 `cd` 到项目目录 (比如 `cd ~/GasOptimization`)，然后直接运行 `code .`。
* **工作方式**: 这样会在 Windows 里打开 VS Code 窗口，但这个窗口是连着 WSL 内部文件系统的。可以直接在图形界面里改代码，然后在 VS Code 自带的终端（或者单独的 Ubuntu 终端）里跑 `forge` 命令，体验很好。

### 3. 网络配置：处理代理问题

安装时最容易卡住的地方就是网络代理。

* **问题原因**: WSL 2 的网络是独立的，它访问不了 Windows 主机上的 `127.0.0.1` 或 `localhost`。
* **解决方法**:
    1.  **找到主机 IP**: 在 Ubuntu 终端里用 `cat /etc/resolv.conf`，`nameserver` 后面的那个 IP 就是 Windows 主机的地址。
    2.  **给 WSL 配代理**: 用找到的 IP 和代理端口（比如 `7078`）设置环境变量：
        ```bash
        export http_proxy="http://<你的主机IP>:7078"
        export https_proxy="http://<你的主机IP>:7078"
        ```
    3.  **配置代理软件**: 这是最关键的一步，需要在 Windows 的代理软件里，把 **“允许来自局域网(LAN)的连接”** 这个选项打开。因为 WSL 的请求在网络上看起来就像是局域网里的另一台设备发过来的。

