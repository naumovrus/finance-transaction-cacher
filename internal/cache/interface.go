package cache

import ent "github.com/naumovrus/finance-transaction-api/internal/entity"

type TransactionCache interface {
	SetTS(key string, value ent.TransactionSend)
	SetTB(key string, value *ent.TransactionTUTO)
	GetTS(key string) *ent.TransactionSend
	GetTB(key string) *ent.TransactionTUTO
	SetCachedData() error
}
