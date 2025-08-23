import * as anchor from "@coral-xyz/anchor";
import { Program } from "@coral-xyz/anchor";
import { SplToken } from "../target/types/spl_token";
import { assert } from "chai";
import { 
  getMint, 
  getOrCreateAssociatedTokenAccount, 
  getAccount 
} from "@solana/spl-token";

describe("spl_token", () => {
  const provider = anchor.AnchorProvider.env();
  anchor.setProvider(provider);
  const program = anchor.workspace.SplToken as Program<SplToken>;
  const mintKeypair = anchor.web3.Keypair.generate();
  const recipient = anchor.web3.Keypair.generate();
  
  // 定义支付账户，以帮助 VS Code 进行类型推断
  const payer = provider.wallet.payer;
  
  let senderAta: anchor.web3.PublicKey;

  it("Is initialized!", async () => {
    const decimals = 9;
    await program.methods
      .initialize(decimals)
      .accounts({
        mint: mintKeypair.publicKey,
        user: provider.wallet.publicKey,
      })
      .signers([mintKeypair])
      .rpc();
    const mintAccount = await getMint(provider.connection, mintKeypair.publicKey);
    assert.equal(mintAccount.mintAuthority!.toBase58(), provider.wallet.publicKey.toBase58());
    assert.equal(mintAccount.decimals, decimals);
    assert.ok(mintAccount.supply === 0n);
  });

  it("Mints tokens!", async () => {
    // **修复**: 使用我们定义的 payer 变量
    const userAtaAccount = await getOrCreateAssociatedTokenAccount(
      provider.connection,
      payer!, // 使用 payer 变量
      mintKeypair.publicKey,
      provider.wallet.publicKey
    );
    senderAta = userAtaAccount.address;

    const mintAmount = new anchor.BN(1000 * Math.pow(10, 9));
    await program.methods
      .mintTokens(mintAmount)
      .accounts({
        mint: mintKeypair.publicKey,
        tokenAccount: senderAta,
        authority: provider.wallet.publicKey,
      })
      .rpc();
    const accountInfo = await getAccount(provider.connection, senderAta);
    assert.ok(accountInfo.amount === BigInt(mintAmount.toString()));
  });

  it("Transfers tokens!", async () => {
    // **修复**: 使用我们定义的 payer 变量
    const recipientAta = await getOrCreateAssociatedTokenAccount(
      provider.connection,
      payer!, // 使用 payer 变量
      mintKeypair.publicKey,
      recipient.publicKey
    );

    const transferAmount = new anchor.BN(500 * Math.pow(10, 9));
    await program.methods
      .transferTokens(transferAmount)
      .accounts({
        fromAccount: senderAta,
        toAccount: recipientAta.address,
        authority: provider.wallet.publicKey,
      })
      .rpc();

    const senderAccountInfo = await getAccount(provider.connection, senderAta);
    const recipientAccountInfo = await getAccount(provider.connection, recipientAta.address);

    assert.ok(senderAccountInfo.amount === BigInt(transferAmount.toString()), "Sender should have 500 tokens left");
    assert.ok(recipientAccountInfo.amount === BigInt(transferAmount.toString()), "Recipient should have 500 tokens");
  });
});
