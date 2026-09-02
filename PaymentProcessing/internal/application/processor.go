package application

import (
	"App/internal/domain"
	"context"
	"fmt"
	"log"
)

const (
	transactionStatusApproved = "approved"
	transactionStatusRejected = "rejected"
)

type Processor struct {
	merchantReader MerchantReader
	merchantWriter MerchantWriter
	txWriter       TransactionWriter
	TxManager TxManager
	cache     Cache[domain.Merchant]
}

func NewProcessor(
	reader MerchantReader,
	writer MerchantWriter,
	txWriter TransactionWriter,
	txManager TxManager,
	cache Cache[domain.Merchant],
) *Processor {
	return &Processor{
		merchantReader: reader,
		merchantWriter: writer,
		txWriter:       txWriter,
		TxManager:      txManager,
		cache:          cache,
	}
}

func (p *Processor) ProcessTransaction(ctx context.Context, tx domain.Transaction) error {
	if p == nil {
		return fmt.Errorf("processor is nil")
	}
	if err := p.validate(); err != nil {
		return err
	}

	merchant, err := p.loadMerchant(ctx, tx.MerchantID)
	if err != nil {
		return err
	}

	if !merchant.IsActive {
		log.Printf("transaction=%d rejected. Merchant %d is not active", tx.ID, tx.MerchantID)
		tx.Status = transactionStatusRejected
		return p.TxManager.WithinTransaction(ctx, func(txCtx context.Context) error {
			return p.txWriter.Save(txCtx, tx)
		})
	}

	newBalance := merchant.Balance + tx.Amount
	tx.Status = transactionStatusApproved

	errUpdateBalance := p.TxManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := p.merchantWriter.UpdateBalance(txCtx, tx.MerchantID, newBalance); err != nil {
			return err
		}
		return p.txWriter.Save(txCtx, tx)
	})
	if errUpdateBalance != nil {
		return errUpdateBalance
	}

	merchant.Balance = newBalance
	p.cache.Put(tx.MerchantID, merchant)
	log.Printf("transaction=%d approved. Merchant %d has been updated. New balance %f", tx.ID, tx.MerchantID, merchant.Balance)

	return nil
}

func (p *Processor) loadMerchant(ctx context.Context, merchantId int64) (domain.Merchant, error) {
	if merchant, ok := p.cache.Get(merchantId); ok {
		return merchant, nil
	}

	merchant, err := p.merchantReader.GetById(ctx, merchantId)
	if err != nil {
		return domain.Merchant{}, fmt.Errorf("load merchant %d: %w", merchantId, err)
	}

	p.cache.Put(merchantId, merchant)
	return merchant, nil
}

func (p *Processor) validate() error {
	switch {
	case p.merchantReader == nil:
		return fmt.Errorf("merchantReader is nil")
	case p.merchantWriter == nil:
		return fmt.Errorf("merchantWriter is nil")
	case p.txWriter == nil:
		return fmt.Errorf("transaction writer is nil")
	case p.TxManager == nil:
		return fmt.Errorf("transaction manager is nil")
	case p.cache == nil:
		return fmt.Errorf("cache is nil")
	default:
		return nil
	}
}
