# 将合约生成GO代码
使用 nodejs，安装 solc 工具：

```go
npm install -g solc
```

使用命令，编译合约代码，会在当目录下生成一个编译好的二进制字节码文件 store_sol_Store.bin：

```go
solcjs --bin Store.sol
```

使用命令，生成合约 abi 文件，会在当目录下生成 store_sol_Store.abi 文件：

```go
solcjs --abi Store.sol
```

abigin 工具可以使用下面的命令安装：

```go
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

使用 abigen 工具根据这两个生成 bin 文件和 abi 文件，生成 go 代码：

```go
abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=store --out=store.go
```

# 部署合约
```go
    auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
    if err != nil {
        log.Fatal(err)
    }
    auth.Nonce = big.NewInt(int64(nonce))
    auth.Value = big.NewInt(0)     // in wei
    auth.GasLimit = uint64(300000) // in units
    auth.GasPrice = gasPrice

    input := "1.0"
    address, tx, instance, err := store.DeployStore(auth, client, input)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(address.Hex())
    fmt.Println(tx.Hash().Hex())
```