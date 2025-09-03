package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/ethclient"
)

func QueryBlockInfoAt() {
	// Read RPC URL from environment (loaded from .env in main)
	url := os.Getenv("ETH_NODE_URL")
	if url == "" {
		log.Fatalf("ETH_NODE_URL is not set. Create a .env file or export the env var before running.")
	}

	// Connect to an Ethereum node
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	blockNumber := big.NewInt(5671744) // Example block number

	// Find info about the block number
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatalf("Failed to retrieve block: %v", err)
	}

	fmt.Println("Block Number: ", block.Number().Uint64())             //5671744
	fmt.Println("Block Hash: ", block.Hash().Hex())                    //0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
	fmt.Println("Block Time: ", block.Time())                          //1527211625
	fmt.Println("Number of Transactions: ", len(block.Transactions())) //144

	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatalf("Failed to retrieve block header: %v", err)
	}

	fmt.Println("Block Number: ", header.Number.Uint64()) //5671744
	fmt.Println("Block Hash: ", header.Hash().Hex())      //0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
	fmt.Println("Block Time: ", header.Time)              //1527211625

	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatalf("Failed to retrieve transaction count: %v", err)
	}

	fmt.Println("Number of Transactions: ", count) //144
}

func transfer() {
	// Read RPC URL from environment (loaded from .env in main)
	url := os.Getenv("ETH_NODE_URL")
	if url == "" {
		log.Fatalf("ETH_NODE_URL is not set. Create a .env file or export the env var before running.")
	}

	// Connect to an Ethereum node
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("Failed to parse PRIVATE_KEY: %v", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to cast public key to ECDSA: %v", err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("From Address:", fromAddress.Hex())
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	toAddress := common.HexToAddress("0xC36512B146C028F651df6ae1Dbd5fB378DB5d583") // Replace with actual receiver address
	var value int64 = 1000000000000000                                             // in wei (0.001 eth)
	var gasLimit uint64 = 21000                                                    // in units

	tx := types.NewTransaction(nonce, toAddress, big.NewInt(value), gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
}
