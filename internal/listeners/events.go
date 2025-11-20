package listeners

import (
	"context"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/amitp-vv/go-backend/internal/models"
	"github.com/amitp-vv/go-backend/internal/services/chain"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Ensures the listener is started only once
var started = false

// StartEventListeners bootstraps live blockchain event listeners
func StartEventListeners() {
	if started {
		return
	}
	started = true

	contract := chain.GetContract()

	// ========== Optional Backfill ==========
	fromBlockEnv := os.Getenv("INDEX_FROM_BLOCK")
	if fromBlockEnv != "" {
		fromBlock := new(big.Int)
		fromBlock.SetString(fromBlockEnv, 10)
		log.Println("Backfilling events from block:", fromBlock.String())

		go backfillEvents(contract, fromBlock)
	}

	// ========== SUBSCRIBE LIVE EVENTS (DISTRIBUTIONCREATED) ==========
	distChan := make(chan *chain.ContractDistributionCreated)
	_, err := contract.WatchDistributionCreated(
		&bind.WatchOpts{Context: context.Background()},
		distChan,
	)
	if err != nil {
		log.Println("Failed to subscribe to DistributionCreated:", err)
		return
	}

	go func() {
		for ev := range distChan {
			handleDistributionCreated(ev)
		}
	}()

	// ========== SUBSCRIBE LIVE EVENTS (CLAIMED) ==========
	claimChan := make(chan *chain.ContractClaimed)
	_, err2 := contract.WatchClaimed(
		&bind.WatchOpts{Context: context.Background()},
		claimChan,
	)
	if err2 != nil {
		log.Println("Failed to subscribe to Claimed:", err2)
		return
	}

	go func() {
		for ev := range claimChan {
			handleClaimed(ev)
		}
	}()

	log.Println("[EVENTS] Blockchain event listeners started.")
}

//
// ==========================================================
// EVENT HANDLERS  (Equivalent to Node.js listeners)
// ==========================================================
//

func handleDistributionCreated(ev *chain.ContractDistributionCreated) {
	id := ev.Id.Uint64()
	snapshotID := ev.SnapshotId.Uint64()
	totalAmount := ev.TotalAmount.String()

	txHash := ""
	if ev.Raw.TxHash != (common.Hash{}) {
		txHash = ev.Raw.TxHash.Hex()
	}

	var dist models.Distribution
	models.DB.Where("onchain_id = ?", id).First(&dist)

	// Create only if not exists
	if dist.ID == 0 {
		dist = models.Distribution{
			PropertyID:  0,
			OnchainID:   id,
			SnapshotID:  snapshotID,
			TotalAmount: totalAmount,
			TxHash:      txHash,
		}
		models.DB.Create(&dist)
		log.Printf("[EVENT] DistributionCreated: saved id=%d\n", id)
	}
}

func handleClaimed(ev *chain.ContractClaimed) {
	id := ev.Id.Uint64()
	user := strings.ToLower(ev.User.Hex())
	amount := ev.Amount.String()

	txHash := ""
	if ev.Raw.TxHash != (common.Hash{}) {
		txHash = ev.Raw.TxHash.Hex()
	}

	claim := models.Claim{
		OnchainID:     id,
		WalletAddress: user,
		Amount:        amount,
		TxHash:        txHash,
	}

	models.DB.Create(&claim)
	log.Printf("[EVENT] Claimed: id=%d user=%s", id, user)
}

//
// ==========================================================
// BACKFILL (Equivalent to backfillEvents() in Node.js)
// ==========================================================
//

func backfillEvents(contract *chain.Contract, fromBlock *big.Int) {
	ctx := context.Background()

	// ---------- DISTRIBUTION CREATED ----------
	distIter, err := contract.FilterDistributionCreated(&bind.FilterOpts{
		Start:   fromBlock.Uint64(),
		Context: ctx,
	})
	if err != nil {
		log.Println("Backfill DistributionCreated failed:", err)
		return
	}

	for distIter.Next() {
		handleDistributionCreated(distIter.Event)
	}

	// ---------- CLAIMED ----------
	claimIter, err2 := contract.FilterClaimed(&bind.FilterOpts{
		Start:   fromBlock.Uint64(),
		Context: ctx,
	})
	if err2 != nil {
		log.Println("Backfill Claimed failed:", err2)
		return
	}

	for claimIter.Next() {
		handleClaimed(claimIter.Event)
	}

	log.Println("[EVENTS] Backfill completed from block:", fromBlock.String())
}
