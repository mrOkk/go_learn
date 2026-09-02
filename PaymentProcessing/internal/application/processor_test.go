package application

import (
	"App/internal/domain"
	"context"
	"errors"
	"strings"
	"testing"

	cachepkg "App/internal/cache"
)

type stubMerchantReader struct {
	merchants map[int64]domain.Merchant
	calls     int
}

func (s *stubMerchantReader) GetById(_ context.Context, id int64) (domain.Merchant, error) {
	s.calls++
	merchant, ok := s.merchants[id]
	if !ok {
		return domain.Merchant{}, errors.New("merchant not found")
	}
	return merchant, nil
}

type stubMerchantWriter struct {
	updates []struct {
		id      int64
		balance float64
	}
	err error
}

func (s *stubMerchantWriter) UpdateBalance(_ context.Context, id int64, balance float64) error {
	s.updates = append(s.updates, struct {
		id      int64
		balance float64
	}{id: id, balance: balance})
	return s.err
}

type stubTransactionWriter struct {
	saved []domain.Transaction
	err   error
}

func (s *stubTransactionWriter) Save(_ context.Context, tx domain.Transaction) error {
	s.saved = append(s.saved, tx)
	return s.err
}

type stubTxManager struct{}

func (s *stubTxManager) WithinTransaction(_ context.Context, fn func(ctx context.Context) error) error {
	return fn(context.Background())
}

func TestProcessorProcessTransactionApproved(t *testing.T) {
	reader := &stubMerchantReader{merchants: map[int64]domain.Merchant{
		1: {Id: 1, Name: "Acme", IsActive: true, Balance: 100},
	}}
	writer := &stubMerchantWriter{}
	txWriter := &stubTransactionWriter{}
	cache := cachepkg.NewLRUCache[domain.Merchant](10)
	processor := NewProcessor(reader, writer, txWriter, &stubTxManager{}, cache)

	tx := domain.Transaction{ID: 42, MerchantID: 1, Amount: 50}
	if err := processor.ProcessTransaction(context.Background(), tx); err != nil {
		t.Fatalf("ProcessTransaction() returned error: %v", err)
	}

	if reader.calls != 1 {
		t.Fatalf("expected 1 merchant read, got %d", reader.calls)
	}

	if len(writer.updates) != 1 {
		t.Fatalf("expected 1 balance update, got %d", len(writer.updates))
	}
	if writer.updates[0].id != 1 || writer.updates[0].balance != 150 {
		t.Fatalf("unexpected balance update: %#v", writer.updates[0])
	}

	if len(txWriter.saved) != 1 {
		t.Fatalf("expected 1 transaction save, got %d", len(txWriter.saved))
	}
	if txWriter.saved[0].Status != "approved" {
		t.Fatalf("expected approved status, got %q", txWriter.saved[0].Status)
	}

	merchant, ok := cache.Get(1)
	if !ok {
		t.Fatal("expected merchant to be stored in cache after approval")
	}
	if merchant.Balance != 150 {
		t.Fatalf("expected cached balance 150, got %f", merchant.Balance)
	}
}

func TestProcessorProcessTransactionRejectedForInactiveMerchant(t *testing.T) {
	reader := &stubMerchantReader{merchants: map[int64]domain.Merchant{
		2: {Id: 2, Name: "Inactive", IsActive: false, Balance: 15},
	}}
	writer := &stubMerchantWriter{}
	txWriter := &stubTransactionWriter{}
	processor := NewProcessor(reader, writer, txWriter, &stubTxManager{}, cachepkg.NewLRUCache[domain.Merchant](10))

	tx := domain.Transaction{ID: 99, MerchantID: 2, Amount: 5}
	if err := processor.ProcessTransaction(context.Background(), tx); err != nil {
		t.Fatalf("ProcessTransaction() returned error: %v", err)
	}

	if len(writer.updates) != 0 {
		t.Fatalf("expected no balance update for inactive merchant, got %d updates", len(writer.updates))
	}

	if len(txWriter.saved) != 1 {
		t.Fatalf("expected 1 transaction save for rejected merchant, got %d", len(txWriter.saved))
	}
	if txWriter.saved[0].Status != "rejected" {
		t.Fatalf("expected rejected status, got %q", txWriter.saved[0].Status)
	}
}

func TestProcessorUsesCacheForMerchantLookup(t *testing.T) {
	reader := &stubMerchantReader{merchants: map[int64]domain.Merchant{
		7: {Id: 7, Name: "Cached", IsActive: true, Balance: 20},
	}}
	writer := &stubMerchantWriter{}
	txWriter := &stubTransactionWriter{}
	cache := cachepkg.NewLRUCache[domain.Merchant](10)
	cache.Put(7, domain.Merchant{Id: 7, Name: "Cached", IsActive: true, Balance: 20})
	processor := NewProcessor(reader, writer, txWriter, &stubTxManager{}, cache)

	tx := domain.Transaction{ID: 11, MerchantID: 7, Amount: 10}
	if err := processor.ProcessTransaction(context.Background(), tx); err != nil {
		t.Fatalf("ProcessTransaction() returned error: %v", err)
	}

	if reader.calls != 0 {
		t.Fatalf("expected cache hit to skip merchantReader.GetById, got %d calls", reader.calls)
	}
}

func TestProcessorValidateNilDependencies(t *testing.T) {
	var processor *Processor
	if err := processor.ProcessTransaction(context.Background(), domain.Transaction{}); err == nil || err.Error() != "processor is nil" {
		t.Fatalf("expected processor is nil error, got %v", err)
	}

	processor = &Processor{}
	err := processor.ProcessTransaction(context.Background(), domain.Transaction{})
	if err == nil || !strings.Contains(err.Error(), "merchantReader is nil") {
		t.Fatalf("expected merchantReader validation error, got %v", err)
	}
}
