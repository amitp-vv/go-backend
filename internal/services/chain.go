package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	contractABI         abi.ABI
	abiOnce             sync.Once
	stableTokenDecimals *uint8
	decimalsMu          sync.Mutex
)

// Load ABI from artifact JSON file
func getAbi() (abi.ABI, error) {
	var err error
	abiOnce.Do(func() {
		path := "../../../SmartContract/artifacts/contracts/RealEstateToken.sol/RealEstateToken.json"
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			err = readErr
			return
		}
		var artifact struct {
			ABI json.RawMessage `json:"abi"`
		}
		if jsonErr := json.Unmarshal(data, &artifact); jsonErr != nil {
			err = jsonErr
			return
		}
		contractABI, err = abi.JSON(bytes.NewReader(artifact.ABI))
	})
	return contractABI, err
}

// Get read-only contract instance
func getContract(client *ethclient.Client, address string) (*bind.BoundContract, error) {
	abi, err := getAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(common.HexToAddress(address), abi, client, client, client), nil
}

// Get signer contract instance
func getSignerContract(client *ethclient.Client, address, pk string) (*bind.BoundContract, *bind.TransactOpts, error) {
	abi, err := getAbi()
	if err != nil {
		return nil, nil, err
	}
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111)) // Sepolia chain ID
	if err != nil {
		return nil, nil, err
	}
	return bind.NewBoundContract(common.HexToAddress(address), abi, client, client, client), auth, nil
}

// Read claimable amount
func ReadClaimableAmount(ctx context.Context, distributionId int64, user string) (string, error) {
	rpc := os.Getenv("SEPOLIA_RPC")
	address := os.Getenv("CONTRACT_ADDRESS")
	if rpc == "" || address == "" {
		return "", errors.New("missing SEPOLIA_RPC or CONTRACT_ADDRESS env")
	}
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return "", err
	}
	contract, err := getContract(client, address)
	if err != nil {
		return "", err
	}
	var out []interface{}
	err = contract.Call(nil, &out, "claimableAmount", distributionId, common.HexToAddress(user))
	if err != nil {
		return "", err
	}
	if len(out) == 0 {
		return "", errors.New("no result")
	}
	return out[0].(*big.Int).String(), nil
}

// Request whitelist update transaction
func RequestWhitelistUpdateTx(ctx context.Context, user string, status bool) (string, error) {
	rpc := os.Getenv("SEPOLIA_RPC")
	pk := os.Getenv("OWNER_PRIVATE_KEY")
	address := os.Getenv("CONTRACT_ADDRESS")
	if rpc == "" || pk == "" || address == "" {
		return "", errors.New("missing env variables")
	}
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return "", err
	}
	contract, auth, err := getSignerContract(client, address, pk)
	if err != nil {
		return "", err
	}
	tx, err := contract.Transact(auth, "requestWhitelistUpdate", common.HexToAddress(user), status)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// Confirm whitelist update transaction
func ConfirmWhitelistUpdateTx(ctx context.Context, user string) (string, error) {
	rpc := os.Getenv("SEPOLIA_RPC")
	pk := os.Getenv("OWNER_PRIVATE_KEY")
	address := os.Getenv("CONTRACT_ADDRESS")
	if rpc == "" || pk == "" || address == "" {
		return "", errors.New("missing env variables")
	}
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return "", err
	}
	contract, auth, err := getSignerContract(client, address, pk)
	if err != nil {
		return "", err
	}
	tx, err := contract.Transact(auth, "confirmWhitelistUpdate", common.HexToAddress(user))
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// Create distribution transaction
func CreateDistributionTx(ctx context.Context, amount string) (string, error) {
	rpc := os.Getenv("SEPOLIA_RPC")
	pk := os.Getenv("OWNER_PRIVATE_KEY")
	address := os.Getenv("CONTRACT_ADDRESS")
	if rpc == "" || pk == "" || address == "" {
		return "", errors.New("missing env variables")
	}
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return "", err
	}
	contract, auth, err := getSignerContract(client, address, pk)
	if err != nil {
		return "", err
	}
	tx, err := contract.Transact(auth, "createDistribution", amount)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// Get stable token decimals (ERC20)
func GetStableTokenDecimals(ctx context.Context) (uint8, error) {
	decimalsMu.Lock()
	defer decimalsMu.Unlock()
	if stableTokenDecimals != nil {
		return *stableTokenDecimals, nil
	}
	rpc := os.Getenv("SEPOLIA_RPC")
	stable := os.Getenv("STABLE_TOKEN_ADDRESS")
	if rpc == "" || stable == "" {
		return 0, errors.New("missing SEPOLIA_RPC or STABLE_TOKEN_ADDRESS env")
	}
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return 0, err
	}
	erc20Abi, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"}]`))
	if err != nil {
		return 0, err
	}
	contract := bind.NewBoundContract(common.HexToAddress(stable), erc20Abi, client, client, client)
	var out []any
	err = contract.Call(nil, &out, "decimals")
	if err != nil {
		return 0, err
	}
	d := out[0].(uint8)
	stableTokenDecimals = &d
	return d, nil
}
