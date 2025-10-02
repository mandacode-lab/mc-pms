package clientmgmt

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/mandacode-ssam/internal/test/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateClientApp_HashesSecretBytesNotEncodedSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	servicePublicID := vo.NewServicePublicID(uuid.New())
	serviceName, _ := vo.NewServiceName("Test Service")
	serviceDesc := "Test Description"

	service := entity.DraftService(serviceName, &serviceDesc)

	secretBytes := []byte("raw-secret-bytes")
	hashedSecret := []byte("hashed-from-raw-bytes")
	encodedSecret := []byte("encoded-secret-for-display")

	// Create mocks
	mockServiceQueryRepo := mock.NewMockServiceQueryRepository(ctrl)
	mockClientAppRepo := mock.NewMockClientAppRepository(ctrl)
	mockClientAppQueryRepo := mock.NewMockClientAppQueryRepository(ctrl)
	mockTxManager := mock.NewMockTransactionManager(ctrl)
	mockSecretGen := mock.NewMockByteRandGen(ctrl)
	mockHasher := mock.NewMockHasher(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)

	usecase := NewUsecase(
		mockClientAppRepo,
		mockClientAppQueryRepo,
		mockServiceQueryRepo,
		mockTxManager,
		mockSecretGen,
		mockHasher,
		mockEncoder,
	)

	req := &in.CreateClientAppRequest{
		ServiceID: servicePublicID,
		Name:      "Test Client",
		Desc:      "Test Description",
	}

	// Set expectations
	mockServiceQueryRepo.EXPECT().FindByPublicID(ctx, servicePublicID).Return(service, nil)
	mockSecretGen.EXPECT().Generate(ctx).Return(secretBytes, nil)

	// CRITICAL: Hash should be called with secretBytes (raw bytes), NOT encodedSecret
	mockHasher.EXPECT().Hash(ctx, secretBytes).Return(hashedSecret, nil)

	// Encode should be called with secretBytes
	mockEncoder.EXPECT().Encode(secretBytes).Return(encodedSecret, nil)

	mockTxManager.EXPECT().WithTx(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(tx out.Tx) error) error {
			return fn(nil)
		},
	)

	mockClientAppRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, tx out.Tx, clientApp *entity.ClientApp) (*entity.ClientApp, error) {
			// Verify that the client app has the hashed secret
			assert.Equal(t, hashedSecret, clientApp.SecretHash())
			return clientApp, nil
		},
	)

	// Execute
	result, err := usecase.CreateClientApp(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, encodedSecret, result.Secret, "Returned secret should be the encoded version")
}
