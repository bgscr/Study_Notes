// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

// ====================================================================
// v5.x 导入路径
// ====================================================================
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";

// 接口和库现在从 'contracts' 主包中导入
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
// ====================================================================

/**
 * @title Stake (质押合约)
 * @dev 一个支持多质押池的系统，并用 MetaNode 代币奖励用户。
 */
contract Stake is Initializable, OwnableUpgradeable, UUPSUpgradeable, PausableUpgradeable {
    using SafeERC20 for IERC20;

    // ====================================================================
    // ## 新增的只读函数 ##
    // ====================================================================
    /**
     * @dev 查看待领取的奖励数量。这是一个只读函数，不消耗 Gas。
     * @param _pid 池的 ID
     * @param _user 用户的地址
     * @return 待领取的奖励数量
     */
    function pendingReward(uint256 _pid, address _user) external view returns (uint256) {
        PoolInfo memory pool = poolInfo[_pid]; // 使用 memory 来避免状态修改
        UserInfo storage user = userInfo[_pid][_user];
        uint256 accMetaNodePerST = pool.accMetaNodePerST;
        uint256 stTokenAmount = pool.stTokenAmount;

        // 如果区块有更新，且池子里有质押，就模拟计算最新的 accMetaNodePerST
        if (block.number > pool.lastRewardBlock && stTokenAmount != 0 && totalPoolWeight != 0) {
            uint256 blockDifference = block.number - pool.lastRewardBlock;
            uint256 totalMetaNodeReward = blockDifference * metaNodePerBlock;
            uint256 poolMetaNodeReward = (totalMetaNodeReward * pool.poolWeight) / totalPoolWeight;
            accMetaNodePerST += (poolMetaNodeReward * 1e18) / stTokenAmount;
        }

        uint256 pending = (user.stAmount * accMetaNodePerST) / 1e18;

        if (pending > user.finishedMetaNode) {
            return pending - user.finishedMetaNode;
        }
        return 0;
    }

    // ... [数据结构和状态变量部分与之前相同，此处省略以保持简洁] ...
    struct UnstakeRequest {
        uint256 amount;
        uint256 unlockBlock;
    }
    struct UserInfo {
        uint256 stAmount;
        uint256 finishedMetaNode;
        uint256 pendingMetaNode;
        UnstakeRequest[] requests;
    }
    struct PoolInfo {
        IERC20 stTokenAddress;
        uint256 poolWeight;
        uint256 lastRewardBlock;
        uint256 accMetaNodePerST;
        uint256 stTokenAmount;
        uint256 minDepositAmount;
        uint256 unstakeLockedBlocks;
    }
    IERC20 public metaNodeToken;
    PoolInfo[] public poolInfo;
    mapping(uint256 => mapping(address => UserInfo)) public userInfo;
    uint256 public totalPoolWeight;
    uint256 public metaNodePerBlock;
    event Deposit(address indexed user, uint256 indexed pid, uint256 amount);
    event Withdraw(address indexed user, uint256 indexed pid, uint256 amount);
    event ClaimReward(address indexed user, uint256 indexed pid, uint256 amount);
    event RequestUnstake(address indexed user, uint256 indexed pid, uint256 amount, uint256 unlockBlock);
    event EmergencyWithdraw(address indexed user, uint256 indexed pid, uint256 amount);


    modifier validatePool(uint256 _pid) {
        // [修正] require 错误信息改回英文
        require(_pid < poolInfo.length, "Stake: pool does not exist");
        _;
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }
    
    function initialize(address _metaNodeToken, uint256 _metaNodePerBlock, address initialOwner) public initializer {
        __Ownable_init(initialOwner);
        __UUPSUpgradeable_init();
        __Pausable_init();
        metaNodeToken = IERC20(_metaNodeToken);
        metaNodePerBlock = _metaNodePerBlock;
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    function poolLength() external view returns (uint256) {
        return poolInfo.length;
    }

    function stake(uint256 _pid, uint256 _amount) external payable whenNotPaused validatePool(_pid) {
        PoolInfo storage pool = poolInfo[_pid];
        UserInfo storage user = userInfo[_pid][msg.sender];

        // [修正] require 错误信息改回英文
        require(_amount >= pool.minDepositAmount, "Stake: amount is less than minimum deposit");

        _updatePool(_pid);
        _updateUser(_pid, msg.sender);

        user.stAmount += _amount;
        pool.stTokenAmount += _amount;

        if (address(pool.stTokenAddress) == address(0)) {
            // [修正] require 错误信息改回英文
            require(msg.value == _amount, "Stake: msg.value must match _amount for native currency");
        } else {
            // [修正] require 错误信息改回英文
            require(msg.value == 0, "Stake: msg.value must be zero for ERC20 tokens");
            pool.stTokenAddress.safeTransferFrom(msg.sender, address(this), _amount);
        }

        emit Deposit(msg.sender, _pid, _amount);
    }

    function unstake(uint256 _pid, uint256 _amount) external whenNotPaused validatePool(_pid) {
        PoolInfo storage pool = poolInfo[_pid];
        UserInfo storage user = userInfo[_pid][msg.sender];

        // [修正] require 错误信息改回英文
        require(user.stAmount >= _amount, "Unstake: amount is greater than staked amount");
        require(_amount > 0, "Unstake: amount must be greater than zero");

        _updatePool(_pid);
        _updateUser(_pid, msg.sender);

        user.stAmount -= _amount;
        pool.stTokenAmount -= _amount;

        uint256 unlockBlock = block.number + pool.unstakeLockedBlocks;
        user.requests.push(UnstakeRequest({amount: _amount, unlockBlock: unlockBlock}));

        emit RequestUnstake(msg.sender, _pid, _amount, unlockBlock);
    }
    
    function withdraw(uint256 _pid) external whenNotPaused validatePool(_pid) {
        PoolInfo storage pool = poolInfo[_pid];
        UserInfo storage user = userInfo[_pid][msg.sender];
        uint256 withdrawableAmount = 0;

        uint256[] memory indicesToClear = new uint256[](user.requests.length);
        uint256 clearCount = 0;

        for (uint i = 0; i < user.requests.length; i++) {
            if (block.number >= user.requests[i].unlockBlock) {
                withdrawableAmount += user.requests[i].amount;
                indicesToClear[clearCount] = i;
                clearCount++;
            }
        }
        
        // [修正] require 错误信息改回英文
        require(withdrawableAmount > 0, "Withdraw: no withdrawable amount");

        for (uint i = clearCount; i > 0; i--) {
            uint index = indicesToClear[i - 1];
            user.requests[index] = user.requests[user.requests.length - 1];
            user.requests.pop();
        }

        if (address(pool.stTokenAddress) == address(0)) {
            payable(msg.sender).transfer(withdrawableAmount);
        } else {
            pool.stTokenAddress.safeTransfer(msg.sender, withdrawableAmount);
        }

        emit Withdraw(msg.sender, _pid, withdrawableAmount);
    }
    
    function claimReward(uint256 _pid) external whenNotPaused validatePool(_pid) {
        _updatePool(_pid);
        _updateUser(_pid, msg.sender);
        
        UserInfo storage user = userInfo[_pid][msg.sender];
        uint256 pending = user.pendingMetaNode;

        // [修正] require 错误信息改回英文
        require(pending > 0, "ClaimReward: no pending rewards");

        user.finishedMetaNode += pending;
        user.pendingMetaNode = 0;
        
        metaNodeToken.safeTransfer(msg.sender, pending);
        
        emit ClaimReward(msg.sender, _pid, pending);
    }
    
    // ... [管理员函数和内部函数与之前相同，此处省略] ...
    function add(address _stTokenAddress, uint256 _poolWeight, uint256 _minDepositAmount, uint256 _unstakeLockedBlocks) external onlyOwner {
        totalPoolWeight += _poolWeight;
        poolInfo.push(
            PoolInfo({
                stTokenAddress: IERC20(_stTokenAddress),
                poolWeight: _poolWeight,
                lastRewardBlock: block.number,
                accMetaNodePerST: 0,
                stTokenAmount: 0,
                minDepositAmount: _minDepositAmount,
                unstakeLockedBlocks: _unstakeLockedBlocks
            })
        );
    }
    function set(uint256 _pid, uint256 _poolWeight, uint256 _minDepositAmount, uint256 _unstakeLockedBlocks) external onlyOwner validatePool(_pid) {
        totalPoolWeight = totalPoolWeight - poolInfo[_pid].poolWeight + _poolWeight;
        poolInfo[_pid].poolWeight = _poolWeight;
        poolInfo[_pid].minDepositAmount = _minDepositAmount;
        poolInfo[_pid].unstakeLockedBlocks = _unstakeLockedBlocks;
    }
    function pause() external onlyOwner {
        _pause();
    }
    function unpause() external onlyOwner {
        _unpause();
    }
    function _updatePool(uint256 _pid) internal {
        PoolInfo storage pool = poolInfo[_pid];
        if (block.number <= pool.lastRewardBlock) {
            return;
        }
        uint256 stTokenAmount = pool.stTokenAmount;
        if (stTokenAmount == 0 || totalPoolWeight == 0) {
            pool.lastRewardBlock = block.number;
            return;
        }
        uint256 blockDifference = block.number - pool.lastRewardBlock;
        uint256 totalMetaNodeReward = blockDifference * metaNodePerBlock;
        uint256 poolMetaNodeReward = (totalMetaNodeReward * pool.poolWeight) / totalPoolWeight;
        pool.accMetaNodePerST += (poolMetaNodeReward * 1e18) / stTokenAmount;
        pool.lastRewardBlock = block.number;
    }
    function _updateUser(uint256 _pid, address _user) internal {
        PoolInfo storage pool = poolInfo[_pid];
        UserInfo storage user = userInfo[_pid][_user];
        uint256 pending = (user.stAmount * pool.accMetaNodePerST) / 1e18;
        if(pending > user.finishedMetaNode) {
            user.pendingMetaNode = pending - user.finishedMetaNode;
        } else {
            user.pendingMetaNode = 0;
        }
    }
}