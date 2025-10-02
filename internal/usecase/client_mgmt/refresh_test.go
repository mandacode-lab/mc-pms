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

func TestRefreshSecret_HashesSecretBytesNotEncodedSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	clientAppID := vo.NewClientAppPublicID(uuid.New())
	serviceID := vo.NewServiceID(1)
	oldHash := []byte("old-hash")

	clientApp := entity.DraftClientApp(serviceID, "Test Client", nil, oldHash)

	secretBytes := []byte("new-raw-secret-bytes")
	hashedSecret := []byte("new-hashed-from-raw-bytes")
	encodedSecret := []byte("new-encoded-secret-for-display")

	// Create mocks
	mockClientAppQueryRepo := mock.NewMockClientAppQueryRepository(ctrl)
	mockClientAppRepo := mock.NewMockClientAppRepository(ctrl)
	mockServiceQueryRepo := mock.NewMockServiceQueryRepository(ctrl)
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

	req := &in.RefreshSecretRequest{
		ClientAppID: clientAppID,
	}

	// Set expectations
	mockClientAppQueryRepo.EXPECT().FindByPublicID(ctx, clientAppID).Return(clientApp, nil)
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

	mockClientAppRepo.EXPECT().Update(ctx, gomock.Any(), clientApp).DoAndReturn(
		func(ctx context.Context, tx out.Tx, app *entity.ClientApp) error {
			// Verify that the secret hash was updated
			assert.Equal(t, hashedSecret, app.SecretHash())
			return nil
		},
	)

	// Execute
	result, err := usecase.RefreshSecret(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, encodedSecret, result.Secret, "Returned secret should be the encoded version")
	assert.Equal(t, hashedSecret, clientApp.SecretHash(), "Client app should have the hashed secret")
}
