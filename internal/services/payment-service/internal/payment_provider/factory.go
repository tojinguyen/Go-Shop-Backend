package paymentprovider

import (
	"fmt"

	"github.com/toji-dev/go-shop/internal/services/payment-service/internal/constant"
)

type PaymentProviderFactory struct {
	providers map[constant.PaymentProviderMethod]PaymentProvider
}

func NewPaymentProviderFactory() *PaymentProviderFactory {
	return &PaymentProviderFactory{
		providers: make(map[constant.PaymentProviderMethod]PaymentProvider),
	}
}

// RegisterProvider đăng ký một provider mới vào factory
func (f *PaymentProviderFactory) RegisterProvider(provider PaymentProvider) {
	f.providers[provider.GetName()] = provider
}

// GetProvider lấy ra một provider dựa trên tên
func (f *PaymentProviderFactory) GetProvider(method constant.PaymentProviderMethod) (PaymentProvider, error) {
	provider, ok := f.providers[method]
	if !ok {
		return nil, fmt.Errorf("payment provider '%s' not found", method)
	}
	return provider, nil
}
