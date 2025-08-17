# 深入解析以太坊三大代币标准：ERC20、ERC721与ERC1155 (附代码学习总结)

## 引言：什么是ERC？

ERC全称为 **Ethereum Request for Comments**（以太坊征求意见稿），是以太坊开发者社区用于提出、讨论和记录新功能、协议或标准的流程。一旦某个ERC被社区广泛接受并最终确定，它就成为了以太坊生态系统中的一个标准。

代币标准（Token Standards）是ERC中最为人熟知的一类，它们定义了一套通用的接口（函数和事件），使得钱包、交易所、去中心化应用（DApps）等能够以一种统一的方式与智能合约进行交互。本文将深入分析和比较三个最核心的代币标准：ERC20、ERC721和ERC1155。

---

## 一、ERC20：同质化代币标准 (Fungible Token)

ERC20是第一个也是应用最广泛的代币标准，它定义了一种**同质化代币**的实现规范。

**核心概念：** 同质化（Fungible）意味着每一个代币单位都是完全相同且可以互换的。例如，你手中的1个USDT和我手中的1个USDT没有任何区别，它们的价值完全相等，可以相互替换。

### 1.1 主要特点
- **可互换性 (Interchangeable)**：任意两个单位的代币价值和功能完全相同。
- **可分割性 (Divisible)**：代币可以被分割成更小的单位（由 `decimals` 字段定义）。
- **统一性 (Uniform)**：所有代幣都共享相同的屬性。

### 1.2 代码接口学习总结 (Interface)

ERC20标准要求智能合约实现以下核心函数和事件：

```solidity
// ERC20 核心接口 (省略了可选的 name, symbol, decimals)

// 函数 (Functions)
function totalSupply() public view returns (uint256);
function balanceOf(address account) public view returns (uint256);
function transfer(address recipient, uint256 amount) public returns (bool);
function allowance(address owner, address spender) public view returns (uint256);
function approve(address spender, uint256 amount) public returns (bool);
function transferFrom(address sender, address recipient, uint256 amount) public returns (bool);

// 事件 (Events)
event Transfer(address indexed from, address indexed to, uint256 value);
event Approval(address indexed owner, address indexed spender, uint256 value);
```

#### 函数解析：

- `totalSupply()`: 返回代币的总供应量。
- `balanceOf(address account)`: 返回指定地址 `account` 的代币余额。
- `transfer(address recipient, uint256 amount)`: 从消息发送者（`msg.sender`）的账户向 `recipient` 地址转移 `amount` 数量的代币。这是最直接的转账方式。
- `approve(address spender, uint256 amount)`: **授权**。允许 `spender` 地址从消息发送者（`msg.sender`）的账户中最多提取 `amount` 数量的代币。这是与DApp交互的核心，例如授权给Uniswap来交易你的代币。
- `allowance(address owner, address spender)`: 查询 `owner` 地址授权给 `spender` 地址的可提取额度。
- `transferFrom(address sender, address recipient, uint256 amount)`: 在获得 `sender` 的授权后，由第三方（`msg.sender`）执行，将 `amount` 数量的代币从 `sender` 账户转移到 `recipient` 账户。这是著名的 **“两步转账”** 模式 (`approve` + `transferFrom`)，是所有DEX和DeFi协议的基础。

### 1.3 应用场景
- **稳定币**: USDT, USDC, DAI
- **治理代币**: UNI (Uniswap), AAVE
- **平台币**: BNB, OKB
- **ICO/IEO发行的项目代币**

### 1.4 优缺点
- **优点**: 简单、标准化程度高，生态系统支持最完善。
- **缺点**:
    - **资产丢失风险**：使用 `transfer` 函数向一个不支持ERC20的合约地址转账，代币会永久锁定在该合约中。
    - **无法代表独特资产**：其同质化的特性使其不适用于表示唯一的、不可互换的资产。

---

## 二、ERC721：非同质化代币标准 (Non-Fungible Token - NFT)

ERC721是为了解决ERC20无法表示独特资产的问题而创建的，它定义了**非同质化代币**（NFT）的规范。

**核心概念：** 非同质化（Non-Fungible）意味着每一个代币都是独一无二、不可分割、不可互换的。每个代币都有一个唯一的 `tokenId` 作为其身份标识。

### 2.1 主要特点
- **唯一性 (Unique)**：每个代币都有一个唯一的 `uint256` 类型的ID (`tokenId`)。
- **不可互换性 (Non-interchangeable)**：`tokenId` 为 1 的代币和 `tokenId` 为 2 的代币是完全不同的资产。
- **所有权 (Ownership)**：记录每个 `tokenId` 的明确所有者。

### 2.2 代码接口学习总结 (Interface)

```solidity
// ERC721 核心接口

// 函数 (Functions)
function balanceOf(address owner) public view returns (uint256 balance);
function ownerOf(uint256 tokenId) public view returns (address owner);
function safeTransferFrom(address from, address to, uint256 tokenId) public;
function transferFrom(address from, address to, uint256 tokenId) public;
function approve(address to, uint256 tokenId) public;
function getApproved(uint256 tokenId) public view returns (address operator);
function setApprovalForAll(address operator, bool _approved) public;
function isApprovedForAll(address owner, address operator) public view returns (bool);

// 事件 (Events)
event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId);
event ApprovalForAll(address indexed owner, address indexed operator, bool approved);
```

#### 函数解析：

- `balanceOf(address owner)`: 返回 `owner` 地址拥有的NFT数量（注意：不是代币的“价值”或“面额”）。
- `ownerOf(uint256 tokenId)`: 返回指定 `tokenId` 的所有者地址。这是ERC721的核心，通过ID追踪所有权。
- `transferFrom(...)`: 与ERC20类似，转移一个NFT的所有权。
- `safeTransferFrom(...)`: 这是ERC721的重要改进。在转账时，它会**检查接收方地址是否为智能合约**。如果是，它会调用该合约的 `onERC721Received` 函数，以确认该合约能够处理NFT，从而避免了ERC20中代币被锁死在不兼容合约中的问题。**强烈推荐使用此函数进行转账**。
- `approve(address to, uint256 tokenId)`: 授权 `to` 地址可以转移指定的 `tokenId`。
- `setApprovalForAll(address operator, bool _approved)`: **批量授权**。授权或取消授权一个 `operator` 地址管理你拥有的**所有**NFT。这在如OpenSea等NFT市场中非常常用，你只需授权一次，市场就可以帮你挂单和转移你的任何NFT。
- `tokenURI(uint256 tokenId)` (在 EIP-721 Metadata Extension 中定义): 返回一个指向JSON文件的URI，该文件描述了NFT的元数据（名称、描述、图片等）。这是NFT显示其艺术价值和属性的关键。

### 2.3 应用场景
- **数字艺术品**: CryptoPunks, Bored Ape Yacht Club (BAYC)
- **游戏道具**: Axie Infinity中的宠物
- **域名服务**: Ethereum Name Service (ENS)
- **身份凭证与证书**: 数字身份、会员资格、门票

### 2.4 优缺点
- **优点**: 完美地代表了数字世界中独特资产的所有权。`safeTransferFrom` 机制更安全。
- **缺点**:
    - **高昂的Gas费**: 每个NFT都是独立的，无论是铸造（mint）还是转移，都需要单独的交易，批量操作非常昂贵。
    - **合约冗余**: 每个NFT系列（Collection）通常需要部署一个独立的智能合约，造成了链上资源的浪费。

---

## 三、ERC1155：多代币标准 (Multi-Token Standard)

ERC1155由Enjin团队提出，旨在解决ERC20和ERC721的局限性，创建一个更高效、更灵活的代币标准。

**核心概念：** ERC1155是一个**多代币标准**，它允许在**单个智能合约**中同时管理无数种同质化（FT）和非同质化（NFT）代币。它引入了 `id` 和 `amount` 的概念，将每个代币类型标识为一个 `id`，并追踪每个地址拥有该 `id` 代币的数量 `amount`。

- 如果一个 `id` 的总供应量 > 1，它就是**同质化代币**（或半同质化代币）。
- 如果一个 `id` 的总供应量 = 1，它就是**非同质化代币**（NFT）。

### 3.1 主要特点
- **高效性 (Efficiency)**：单个合约管理多种代币，极大减少了部署成本和链上数据冗余。
- **批量操作 (Batch Operations)**：支持一次性转移多种、多个代币，显著降低了Gas费用。
- **混合性 (Hybrid)**：天然支持同质化、非同质化及半同质化代币。
- **安全性 (Safety)**：强制要求使用 `safeTransfer` 类型的函数，并引入接收钩子（receiver hook）来防止资产丢失。

### 3.2 代码接口学习总结 (Interface)

```solidity
// ERC1155 核心接口

// 函数 (Functions)
function balanceOf(address account, uint256 id) public view returns (uint256);
function balanceOfBatch(address[] memory accounts, uint256[] memory ids) public view returns (uint256[] memory);
function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes memory data) public;
function safeBatchTransferFrom(address from, address to, uint256[] memory ids, uint256[] memory amounts, bytes memory data) public;
function setApprovalForAll(address operator, bool approved) public;
function isApprovedForAll(address account, address operator) public view returns (bool);

// 事件 (Events)
event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value);
event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values);
event ApprovalForAll(address indexed account, address indexed operator, bool approved);
event URI(string value, uint256 indexed id);
```

#### 函数解析：

- `balanceOf(address account, uint256 id)`: 查询 `account` 地址拥有的 `id` 类型代币的数量。这是与ERC20/ERC721最显著的区别，查询时必须同时提供地址和代币ID。
- `balanceOfBatch(...)`: **批量查询**多个地址的多种代币余额，非常高效。
- `safeTransferFrom(..., uint256 id, uint256 amount, ...)`: 转移**单一类型**的代币。`id` 指定代币类型，`amount` 指定数量。
- `safeBatchTransferFrom(..., uint256[] memory ids, uint256[] memory amounts, ...)`: **批量转移**。ERC1155的“杀手级”功能，允许在单次交易中转移多种不同类型的代币给一个接收者，极大地节省了Gas。
- `setApprovalForAll(...)`: 与ERC721类似，授权一个 `operator` 管理你的**所有**代币（所有 `id`）。注意，ERC1155没有像ERC20那样的单笔 `approve`，因为批量操作更为常用。
- **接收钩子**: 类似于 `safeTransferFrom`，当向合约转账时，会调用接收方合约的 `onERC1155Received` 或 `onERC1155BatchReceived` 函数来确认其能处理这些代币。

### 3.3 应用场景
- **区块链游戏**: 这是ERC1155最主要的应用领域。一个游戏合约可以同时管理：
    - ID 1: 金币 (同质化, FT)
    - ID 2: 魔法药水 (同质化, FT)
    - ID 1001: 传奇宝剑 (非同质化, NFT)
    - ID 1002: 史诗盔甲 (非同质化, NFT)
- **票务系统**: 演唱会门票，其中普通票是FT，而VIP座位票可以是NFT。
- **数字组合资产**: 将多种代币打包成一个资产组合进行交易。

### 3.4 优缺点
- **优点**: 极高的Gas效率，灵活性强，一个合约即可支持复杂经济系统。
- **缺点**:
    - **复杂性更高**: 实现和理解起来比前两者更复杂。
    - **所有权追踪**: 追踪单个NFT（供应量为1的`id`）的所有者不如ERC721直观，需要额外的数据结构或事件索引。

---

## 四、三大协议对比总结

| 特性 / 协议 | ERC20 (同质化代币) | ERC721 (非同质化代币) | ERC1155 (多代币) |
| :--- | :--- | :--- | :--- |
| **代币类型** | 同质化 (Fungible) | 非同质化 (Non-Fungible) | 同质化 & 非同质化 & 半同质化 |
| **核心标识** | 无独立标识，只关心数量 | 唯一的 `tokenId` | `id` (代币类型) + `amount` (数量) |
| **主要应用** | 货币、稳定币、治理代币 | 数字艺术品、收藏品、域名 | 游戏道具、票务、组合资产 |
| **转账安全** | `transfer` 有风险 | `safeTransferFrom` (推荐) | `safeTransferFrom` (强制) |
| **授权机制** | `approve` (单笔) | `approve` (单个tokenId) / `setApprovalForAll` (全部) | `setApprovalForAll` (全部) |
| **批量操作** | 不支持 | 不支持 | **原生支持 (核心优势)** |
| **Gas 效率** | 中等 | 低 (单笔操作昂贵) | **高 (批量操作极省Gas)** |
| **合约模型** | 1个合约 : 1种代币 | 1个合约 : 1个NFT系列 | **1个合约 : N种代币** |

## 最终代码学习总结与选择建议

- **从演进角度看**：`ERC20 -> ERC721 -> ERC1155` 是一个功能不断增强、效率不断优化的演进路径。ERC721解决了ERC20无法表达唯一性的问题，而ERC1155则通过批量处理和多代币合一，解决了ERC20和ERC721在特定场景（尤其是游戏）下的低效率问题。

- **从核心代码模式看**：
    - **ERC20** 的核心是 `balanceOf(owner)` 和 `approve` + `transferFrom` 的两步授权模式。
    - **ERC721** 的核心是 `ownerOf(tokenId)` 来追踪所有权，以及 `tokenURI` 来链接元数据。`safeTransferFrom` 是安全实践的关键。
    - **ERC1155** 的核心是通过 `(id, amount)` 的元组来定义资产，并通过 `safeBatchTransferFrom` 实现高效的原子化批量操作。

- **如何选择？**
    - 如果你正在构建一个**货币系统、稳定币或治理代币**，选择 **ERC20**。它的简单性和广泛的生态支持是无与伦比的。
    - 如果你正在创建一个**数字艺术品、收藏品、或者任何需要证明唯一所有权**的 DApp，选择 **ERC721**。它在表达独特性方面是黄金标准。
    - 如果你正在开发一个**区块链游戏、复杂的DeFi协议或任何需要在一个合约中处理多种类型资产**的系统，毫无疑问应该选择 **ERC1155**。它的高效性和灵活性将为你节省大量的开发成本和用户Gas费用。