// 导入 anchor-lang 库，这是 Anchor 框架的核心
use anchor_lang::prelude::*;
// 导入 anchor-spl 库中的必要模块
use anchor_spl::token::{self, Mint, Token, TokenAccount};

// 你的程序 ID
declare_id!("8AgJackVZieemuZKzKuHF8Li8fBCXxs4Q2y4EyTqBcRf");

#[program]
pub mod spl_token {
    use super::*;

    // `initialize` 函数保持不变
    pub fn initialize(ctx: Context<Initialize>, _decimals: u8) -> Result<()> {
        msg!("Token mint initialized successfully by Anchor!");
        msg!("New Mint Account: {}", ctx.accounts.mint.key());
        Ok(())
    }

    // `mint_tokens` 函数保持不变
    pub fn mint_tokens(ctx: Context<MintTokens>, amount: u64) -> Result<()> {
        msg!("Minting {} tokens...", amount);
        let cpi_accounts = token::MintTo {
            mint: ctx.accounts.mint.to_account_info(),
            to: ctx.accounts.token_account.to_account_info(),
            authority: ctx.accounts.authority.to_account_info(),
        };
        let cpi_context = CpiContext::new(ctx.accounts.token_program.to_account_info(), cpi_accounts);
        token::mint_to(cpi_context, amount)?;
        msg!("Tokens minted successfully.");
        Ok(())
    }

    // **新增**: `transfer_tokens` 函数
    pub fn transfer_tokens(ctx: Context<TransferTokens>, amount: u64) -> Result<()> {
        msg!("Transferring {} tokens...", amount);

        // 创建 CPI 调用来请求 Token Program 为我们转移代币
        let cpi_accounts = token::Transfer {
            from: ctx.accounts.from_account.to_account_info(),
            to: ctx.accounts.to_account.to_account_info(),
            authority: ctx.accounts.authority.to_account_info(),
        };
        let cpi_context = CpiContext::new(ctx.accounts.token_program.to_account_info(), cpi_accounts);
        
        // 执行 CPI 调用
        token::transfer(cpi_context, amount)?;

        msg!("Tokens transferred successfully.");
        Ok(())
    }
}

// `Initialize` 结构体保持不变
#[derive(Accounts)]
#[instruction(decimals: u8)]
pub struct Initialize<'info> {
    #[account(init, payer = user, mint::authority = user, mint::decimals = decimals)]
    pub mint: Account<'info, Mint>,
    #[account(mut)]
    pub user: Signer<'info>,
    pub system_program: Program<'info, System>,
    pub token_program: Program<'info, Token>,
    pub rent: Sysvar<'info, Rent>,
}

// `MintTokens` 结构体保持不变
#[derive(Accounts)]
pub struct MintTokens<'info> {
    #[account(mut)]
    pub mint: Account<'info, Mint>,
    pub authority: Signer<'info>,
    #[account(mut)]
    pub token_account: Account<'info, TokenAccount>,
    pub token_program: Program<'info, Token>,
}

// **新增**: `TransferTokens` 账户结构体
#[derive(Accounts)]
pub struct TransferTokens<'info> {
    // 谁是这笔交易的发起者和签名者？
    pub authority: Signer<'info>,

    // 从哪个代币账户转出代币
    // `mut` 表示它的余额会减少
    #[account(mut)]
    pub from_account: Account<'info, TokenAccount>,

    // 转入到哪个代币账户
    // `mut` 表示它的余额会增加
    #[account(mut)]
    pub to_account: Account<'info, TokenAccount>,
    
    pub token_program: Program<'info, Token>,
}
