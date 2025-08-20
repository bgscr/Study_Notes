// SPDX-License-Identifier: MIT
// 指定软件许可，MIT表示代码开源且限制宽松
pragma solidity ^0.8.20;
// 指定Solidity编译器版本，^表示兼容0.8.20及以上但不包括0.9.0的版本

// --- 导入外部合约 ---
// 从OpenZeppelin库导入标准的ERC20代币合约，提供了基础的代币功能
import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.9.3/contracts/token/ERC20/ERC20.sol";
// 导入Ownable合约，用于实现合约所有权管理，方便设置管理员权限
import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.9.3/contracts/access/Ownable.sol";


// --- 接口定义 ---
// 定义Uniswap V2工厂合约的接口，只需要createPair函数来创建交易对
interface IUniswapV2Factory {
    function createPair(address tokenA, address tokenB) external returns (address pair);
}

// 定义Uniswap V2路由合约的接口，主要用于添加流动性
interface IUniswapV2Router02 {
    function factory() external pure returns (address);
    function WETH() external pure returns (address);
    function addLiquidityETH(
        address token,
        uint amountTokenDesired,
        uint amountTokenMin,
        uint amountETHMin,
        address to,
        uint deadline
    ) external payable returns (uint amountToken, uint amountETH, uint liquidity);
}

// --- 主合约定义 ---
// 定义我们的Meme代币合约，它继承自ERC20和Ownable合约
contract MyShibToken is ERC20, Ownable {

    // --- 状态变量 ---

    // Uniswap V2路由合约的实例。
    // 主网地址: 0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D
    // Sepolia测试网地址: 0xC532a74256D3Db42D0Bf7a0400fEFDbad7694008
    IUniswapV2Router02 public uniswapV2Router;
    // 本代币与WETH在Uniswap上的交易对地址
    address public uniswapV2Pair;

    // 一个标志位，用于控制是否开启交易。可用于项目启动时的控制。
    bool private _isTradingEnabled;

    // 单笔最大交易量限制。设置为0则禁用此功能。
    uint256 public maxTxAmount;
    // 单个钱包最大持币量限制。设置为0则禁用此功能。
    uint256 public maxWalletAmount;

    // --- 税率设置 ---
    // 税率单位为基点 (100 = 1%)
    uint256 public liquidityFee = 200; // 2% 流动性税
    uint256 public marketingFee = 200; // 2% 营销税
    uint256 public redistributionFee = 100; // 1% 分红税 (实现较为复杂，此处作为占位符)
    uint256 public totalFees = liquidityFee + marketingFee + redistributionFee; // 总税率
    
    // 接收营销税费的钱包地址
    address public marketingWallet;

    // 一个映射，用于将某些地址排除在税费计算之外（例如：合约所有者、合约自身、交易对地址）
    mapping(address => bool) private _isExcludedFromFee;

    // 一个映射，用于将某些地址排除在交易限制之外（例如：交易对地址）
    mapping(address => bool) private _isExcludedFromLimits;

    // --- 事件定义 ---
    event TradingEnabled(bool enabled); // 交易状态变更事件
    event FeesUpdated(uint256 liquidity, uint256 marketing, uint256 redistribution); // 税率更新事件
    event MaxTxUpdated(uint256 amount); // 最大交易量更新事件

    // --- 构造函数 ---
    // 在合约部署时执行一次，用于初始化
    constructor(
        string memory _name,          // 代币名称
        string memory _symbol,        // 代币符号
        uint256 _initialSupply,   // 初始供应量
        address _routerAddress,   // Uniswap路由合约地址
        address _marketingWallet  // 营销钱包地址
    ) ERC20(_name, _symbol) Ownable() { // 初始化继承的合约
        // 1. 将初始供应量的代币铸造给合约自身，以便后续添加流动性
        _mint(address(this), _initialSupply * (10**decimals()));
        
        // 2. 设置Uniswap路由合约地址
        uniswapV2Router = IUniswapV2Router02(_routerAddress);

        // 3. 设置营销钱包地址
        require(_marketingWallet != address(0), "Marketing wallet address cannot be the zero address");
        marketingWallet = _marketingWallet;

        // 4. 设置交易和钱包持币量限制
        uint256 _totalSupply = totalSupply();
        maxTxAmount = (_totalSupply * 5) / 100; // 单笔最大交易量为总量的5% (已从2%上调)
        maxWalletAmount = (_totalSupply * 5) / 100; // 单钱包最大持币量为总量的5%

        // 5. 将合约所有者、合约自身和营销钱包地址排除在税费和交易限制之外
        _isExcludedFromFee[owner()] = true;
        _isExcludedFromFee[address(this)] = true;
        _isExcludedFromFee[marketingWallet] = true;

        _isExcludedFromLimits[owner()] = true;
        _isExcludedFromLimits[address(this)] = true;
        _isExcludedFromLimits[marketingWallet] = true;

        // 6. 初始状态下，交易是关闭的
        _isTradingEnabled = false; 
    }

    // --- 重写ERC20的_transfer函数 ---
    // 这是所有代币转移的核心逻辑，我们在这里加入税费和限制功能
    function _transfer(
        address from,
        address to,
        uint256 amount
    ) internal override {
        require(from != address(0), "ERC20: transfer from the zero address");
        require(to != address(0), "ERC20: transfer to the zero address");
        require(amount > 0, "Transfer amount must be greater than zero");

        // --- 交易开关检查 ---
        if (!_isTradingEnabled) {
            // 允许合约所有者和合约自身在交易关闭时进行操作（例如：添加初始流动性）
            require(_isExcludedFromFee[from] || _isExcludedFromFee[to], "Trading is not yet enabled");
        }
        
        // --- 交易量和钱包持仓限制检查 ---
        if (!_isExcludedFromLimits[from] && !_isExcludedFromLimits[to]) {
            require(amount <= maxTxAmount, "Transfer amount exceeds the max transaction amount");
            // 检查买入或普通转账，卖出到交易对地址时不检查目标钱包持仓
            if (to != uniswapV2Pair) { 
                require(balanceOf(to) + amount <= maxWalletAmount, "Exceeds maximum wallet amount");
            }
        }

        // --- 税费逻辑 ---
        bool takeFee = true;
        // 如果发送方或接收方在免税名单中，则不收取费用
        if(_isExcludedFromFee[from] || _isExcludedFromFee[to]) {
            takeFee = false;
        }

        if (takeFee) {
            // 计算需要收取的总税费
            uint256 fees = (amount * totalFees) / 10000;
            // 计算扣除税费后，接收方实际收到的金额
            uint256 taxedAmount = amount - fees;

            // 将税费部分发送到本合约地址暂存
            super._transfer(from, address(this), fees);
            // 将扣税后的金额发送给接收方
            super._transfer(from, to, taxedAmount);
            
            // 处理合约中暂存的税费
            _handleFees(fees);

        } else {
            // 如果无需收费，则执行标准的转账操作
            super._transfer(from, to, amount);
        }
    }

    // --- 内部函数：处理税费分配 ---
    function _handleFees(uint256 fees) private {
        // 这是一个简化的税费处理机制。
        if (totalFees == 0) return; // 避免除以零
        
        uint256 marketingShare = (fees * marketingFee) / totalFees; // 计算营销部分的份额

        // 将营销份额的代币发送到营销钱包
        if (marketingShare > 0) {
            super._transfer(address(this), marketingWallet, marketingShare);
        }
    }


    // --- 仅限所有者调用的函数 ---

    /**
     * @notice [核心启动函数] 创建流动性池并开启交易。
     * @param tokenAmount 用于添加流动性的代币数量。
     */
    function setupLiquidityAndEnableTrading(uint256 tokenAmount) external payable onlyOwner {
        require(!_isTradingEnabled, "Trading is already enabled");
        require(uniswapV2Pair == address(0), "Liquidity already added");

        // 1. 创建Uniswap交易对
        uniswapV2Pair = IUniswapV2Factory(uniswapV2Router.factory())
            .createPair(address(this), uniswapV2Router.WETH());

        // 2. 将交易对地址排除在交易限制和税费之外
        _isExcludedFromLimits[uniswapV2Pair] = true;
        _isExcludedFromFee[uniswapV2Pair] = true;

        // 3. 授权路由合约可以从本合约地址转移指定数量的代币
        _approve(address(this), address(uniswapV2Router), tokenAmount);

        // 4. 调用路由合约的addLiquidityETH函数添加流动性
        uniswapV2Router.addLiquidityETH{value: msg.value}(
            address(this),
            tokenAmount,
            0, // 接受任何滑点，因为这是初始流动性
            0, // 接受任何滑点
            owner(), // LP代币会发送给合约所有者
            block.timestamp
        );

        // 5. 开启交易
        _isTradingEnabled = true;
        emit TradingEnabled(true);
    }

    /**
     * @notice 更新税率百分比。
     * @param _liquidity 新的流动性税率（基点）。
     * @param _marketing 新的营销税率（基点）。
     * @param _redistribution 新的分红税率（基点）。
     */
    function setFees(uint256 _liquidity, uint256 _marketing, uint256 _redistribution) external onlyOwner {
        liquidityFee = _liquidity;
        marketingFee = _marketing;
        redistributionFee = _redistribution;
        totalFees = _liquidity + _marketing + _redistribution;
        require(totalFees <= 1000, "Total fees cannot exceed 10%"); // 设置一个上限，防止设置过高税率
        emit FeesUpdated(_liquidity, _marketing, _redistribution); // 触发税率更新事件
    }
    
    /**
     * @notice 更新单笔最大交易量。
     * @param _amount 新的单笔最大交易量。
     */
    function setMaxTxAmount(uint256 _amount) external onlyOwner {
        // 确保最大交易量不会设置得过低
        require(_amount >= (totalSupply() * 1) / 1000, "Max Tx Amount cannot be less than 0.1% of total supply");
        maxTxAmount = _amount;
        emit MaxTxUpdated(_amount); // 触发最大交易量更新事件
    }

    /**
     * @notice 将一个账户设置为免税或取消免税。
     */
    function excludeFromFee(address account, bool excluded) external onlyOwner {
        _isExcludedFromFee[account] = excluded;
    }

    /**
     * @notice 将一个账户设置为不受交易限制或取消。
     */
    function excludeFromLimits(address account, bool excluded) external onlyOwner {
        _isExcludedFromLimits[account] = excluded;
    }

    /**
     * @notice 提取合约中剩余的代币到所有者钱包。
     */
    function withdrawRemainingTokens() external onlyOwner {
        uint256 remainingBalance = balanceOf(address(this));
        require(remainingBalance > 0, "No tokens left in contract");
        super._transfer(address(this), owner(), remainingBalance);
    }

    // --- 回退函数 ---
    // 允许合约直接接收ETH（例如，用于未来可能的兑换操作）
    receive() external payable {}
}
