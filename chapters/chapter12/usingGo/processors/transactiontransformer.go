package processors

import (
	"github.com/karpov-kir/kafka-in-action/models"
	"github.com/karpov-kir/kafka-in-action/utils"
)

type transactionTransformer struct {
	getFunds   func() models.Funds
	storeFunds func(models.Funds)
}

func newTransactionTransformer(getFunds func() models.Funds, storeFunds func(models.Funds)) *transactionTransformer {
	return &transactionTransformer{
		getFunds:   getFunds,
		storeFunds: storeFunds,
	}
}

func (tt *transactionTransformer) transformTransaction(transaction models.Transaction) models.TransactionResult {
	if transaction.Type == models.TransactionTypeDEPOSIT {
		return models.TransactionResult{
			Transaction: transaction,
			Funds:       tt.depositFunds(transaction),
			Success:     true,
		}
	}

	if tt.hasEnoughFunds(transaction) {
		return models.TransactionResult{
			Transaction: transaction,
			Funds:       tt.withdrawFunds(transaction),
			Success:     true,
		}
	}

	unionNullErrorType := models.NewUnionNullErrorType()
	unionNullErrorType.UnionType = models.UnionNullErrorTypeTypeEnumErrorType
	unionNullErrorType.ErrorType = models.ErrorTypeINSUFFICIENT_FUNDS

	return models.TransactionResult{
		Transaction: transaction,
		Funds:       tt.getFunds(),
		Success:     false,
		ErrorType:   unionNullErrorType,
	}
}

func (tt *transactionTransformer) depositFunds(transaction models.Transaction) models.Funds {
	newBalance := utils.DecimalToAvroBytes(utils.AvroBytesToDecimal(tt.getFunds().Balance) + utils.AvroBytesToDecimal(transaction.Amount))
	newFunds := models.Funds{
		Balance: newBalance,
		Account: transaction.Account,
	}
	tt.storeFunds(newFunds)
	return newFunds
}

func (tt *transactionTransformer) withdrawFunds(transaction models.Transaction) models.Funds {
	newBalance := utils.DecimalToAvroBytes(utils.AvroBytesToDecimal(tt.getFunds().Balance) - utils.AvroBytesToDecimal(transaction.Amount))
	newFunds := models.Funds{
		Balance: newBalance,
		Account: transaction.Account,
	}
	tt.storeFunds(newFunds)
	return newFunds
}

func (tt *transactionTransformer) hasEnoughFunds(transaction models.Transaction) bool {
	return utils.AvroBytesToDecimal(tt.getFunds().Balance) >= utils.AvroBytesToDecimal(transaction.Amount)
}
