# **Substrate 自定义区块链项目文档**

本文档详细介绍了基于 polkadot-sdk-solochain-template 构建的自定义区块链项目。该项目在标准节点模板的基础上，实现了一个自定义的“资产登记”功能模块 (Pallet)。

### **一、 项目架构**

本项目遵循 Substrate FRAME 框架的标准架构，该架构将区块链节点的核心组件解耦，实现了高度的模块化。

1. **节点 (Node)**:  
   * **位置**: node/ 目录  
   * **作用**: 节点是运行区块链网络的可执行程序。它包含了所有底层逻辑，如网络通信 (P2P)、交易池、共识引擎 (BABE 和 GRANDPA) 以及执行环境。它负责接收交易、打包区块并与其他节点同步状态。  
2. **运行时 (Runtime)**:  
   * **位置**: runtime/ 目录  
   * **作用**: 运行时是区块链的“状态转换函数”，定义了链的核心业务逻辑。它由多个 Pallet 组合而成，并被编译成 WASM 二进制文件存储在链上，这使得区块链可以实现无分叉的平滑升级。  
3. **Pallets (模块)**:  
   * **位置**: pallets/ 目录  
   * **作用**: Pallet 是 Substrate 中最小的功能单元。每个 Pallet 封装了一组特定的存储项、交易和事件。我们的自定义功能就是通过修改 pallets/template 实现的。

### **二、 功能说明**

本项目的核心是在标准节点的基础上，通过 pallet-template 实现了一个简易的**链上资产登记和转移系统**。

* **核心功能**:  
  1. **资产注册 (register\_asset)**: 允许任何账户支付交易费，将一个唯一的数字 ID (u32) 注册为资产，并将自己记录为该资产的所有者。  
  2. **资产转移 (transfer\_asset)**: 允许已注册资产的当前所有者，将该资产的所有权转移给另一个指定的账户。  
* **业务逻辑与约束**:  
  * 资产 ID 必须是唯一的，无法重复注册。  
  * 只有资产的当前所有者才有权发起对该资产的转移操作。  
  * 所有成功的注册和转移操作都会在链上发出对应的事件 (AssetRegistered, AssetTransferred)，以便链下应用可以监听和响应。

### **三、 代码结构**

项目的关键代码位于以下文件中：

1. **自定义 Pallet 核心逻辑**:  
   * **文件**: pallets/template/src/lib.rs  
   * **内容**: 这是我们实现所有自定义功能的地方。  
     * \#\[pallet::storage\] 部分的 AssetOwner 定义了资产 ID 到所有者的链上存储映射。  
     * \#\[pallet::event\] 部分定义了 AssetRegistered 和 AssetTransferred 事件。  
     * \#\[pallet::error\] 部分定义了 AssetIdInUse, NotTheOwner 等业务逻辑错误。  
     * \#\[pallet::call\] 部分实现了 register\_asset 和 transfer\_asset 这两个可供用户调用的交易函数。  
2. **运行时 Pallet 集成**:  
   * **文件**: runtime/src/lib.rs  
   * **内容**: 这是“主板”文件，它使用 construct\_runtime\! 宏将所有 Pallet（包括我们修改过的 TemplateModule）组装成一个完整的运行时。  
3. **运行时配置**:  
   * **文件**: runtime/src/configs/mod.rs  
   * **内容**: 在这个文件中，我们为 pallet-template 实现了 Config trait，将其集成到 Runtime 中，并配置了它所需的 RuntimeEvent 类型。

### **四、 部署和运行步骤**

1. **环境准备**:  
   * **安装核心依赖**: 运行官方脚本来安装 Rust 和其他基础工具。  
     curl https://getsubstrate.io \-sSf | bash \-s \-- \--fast

   * **安装 WASM 编译目标**: 这是编译 Substrate 运行时所必需的。  
     rustup target add wasm32-unknown-unknown

   * **安装 Rust 源码**: 编译 WASM 同样需要 Rust 的标准库源码。  
     rustup component add rust-src

2. **获取代码**:  
   * 使用 Git 克隆最新的官方 Solochain 模板：  
     git clone https://github.com/paritytech/polkadot-sdk-solochain-template.git solochain-template  
     cd solochain-template

3. **编译节点**:  
   * 在项目根目录下，执行编译命令。首次编译会耗时较长。（如果遇到网络问题，请参考下一章节的解决方案）  
     cargo build \--release

4. **启动本地开发节点**:  
   * 使用以下命令启动一个单节点的本地开发链。--dev 参数会清除旧数据，--rpc-external 允许外部 UI 连接。  
     ./target/release/solochain-template-node \--dev \--rpc-external

   * 保持此终端窗口运行。  
5. **与节点交互**:  
   * 在 WSL 环境中安装并启动 Firefox 浏览器 (sudo apt install firefox 然后 firefox)。  
   * 在 Firefox 中访问 Polkadot-JS Apps: [https://polkadot.js.org/apps/](https://polkadot.js.org/apps/)。  
   * 在页面左上角连接到 Development \> Local Node。  
   * 使用 Developer \> Extrinsics 页面来提交 template.registerAsset 和 template.transferAsset 交易。  
   * 使用 Developer \> Chain state 页面来查询 template.assetOwner 存储项，验证链上数据的变化。

### **五、 常见问题与解决方案 (Troubleshooting)**

在环境配置和编译过程中，可能会遇到一些问题，以下是本次实践中遇到的问题及其解决方案。

1. **问题**: cargo build 下载依赖时因网络问题超时或失败。  
   * **解决方案**: 配置 cargo 使用国内的镜像源。  
     1. 创建或编辑配置文件 \~/.cargo/config.toml。  
     2. 在文件中添加以下内容，以使用清华大学的 TUNA 镜像源为例：  
        \[source.crates-io\]  
        replace-with \= "tuna"

        \[source.tuna\]  
        registry \= "https://mirrors.tuna.tsinghua.edu.cn/git/crates.io-index.git"

2. **问题**: getsubstrate.io 安装脚本报错 E: Unable to locate package protobuf。  
   * **解决方案**: 手动安装 protobuf-compiler，并通过 equivs 创建一个虚拟包来满足脚本的依赖检查。  
     \# 1\. 安装真正的编译器  
     sudo apt-get update && sudo apt-get install \-y protobuf-compiler

     \# 2\. 安装虚拟包创建工具  
     sudo apt install equivs

     \# 3\. 创建虚拟包配置文件  
     equivs-control protobuf

     \# 4\. (可选) 编辑 protobuf 文件，通常默认即可

     \# 5\. 构建并安装虚拟包  
     equivs-build protobuf  
     sudo dpkg \-i protobuf\_\*.deb

   * **注意**: 完成这些步骤后，再重新运行 getsubstrate.io 脚本。  
3. **问题**: 在 Windows 浏览器中打开 Polkadot-JS Apps，无法连接到在 WSL 中运行的节点，一直显示 "Initializing connection"。  
   * **根本原因**: Windows 和 WSL 属于不同的网络环境，直接访问 127.0.0.1 无法找到对方。  
   * **最佳解决方案**: 在 WSL 内部运行一个图形界面的浏览器。  
     1. **安装 Firefox**:  
        sudo apt update && sudo apt install firefox

     2. **启动 Firefox**:  
        firefox

     3. 在弹出的 Firefox 窗口中访问 Polkadot-JS Apps 网站，并连接到默认的 Local Node (ws://127.0.0.1:9944) 即可成功。