package clientaccess

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/test/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVerifyClient_ComparesWithRawSecretBytes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	clientAppID := vo.NewClientAppPublicID(uuid.New())
	serviceID := vo.NewServiceID(1)

	// Simulate the flow:
	// 1. During creation: secretBytes -> Hash(secretBytes) -> hashedSecret (stored in DB)
	// 2. During creation: secretBytes -> Encode(secretBytes) -> encodedSecret (returned to user)
	// 3. During verify: user sends secretBytes in Basic Auth
	// 4. Compare(hashedSecret, secretBytes) should succeed
	secretBytes := []byte("raw-secret-bytes")
	hashedSecret := []byte("hashed-from-raw-bytes")

	clientApp := entity.DraftClientApp(serviceID, "Test Client", nil, hashedSecret)

	serviceName, _ := vo.NewServiceName("Test Service")
	service := entity.DraftService(serviceName, nil)
	servicePublicID := service.PublicID()

	// Create mocks
	mockClientAppQueryRepo := mock.NewMockClientAppQueryRepository(ctrl)
	mockServiceQueryRepo := mock.NewMockServiceQueryRepository(ctrl)
	mockHasher := mock.NewMockHasher(ctrl)

	usecase := NewUsecase(
		mockClientAppQueryRepo,
		mockServiceQueryRepo,
		mockHasher,
	)

	req := &in.VerifyClientRequest{
		ClientID:     clientAppID,
		ClientSecret: secretBytes, // User sends raw bytes (not encoded)
	}

	// Set expectations
	mockClientAppQueryRepo.EXPECT().FindByPublicID(ctx, clientAppID).Return(clientApp, nil)

	// CRITICAL: Compare should be called with hashedSecret and raw secretBytes
	mockHasher.EXPECT().Compare(ctx, hashedSecret, secretBytes).Return(nil)

	mockServiceQueryRepo.EXPECT().FindByID(ctx, serviceID).Return(service, nil)

	// Execute
	result, err := usecase.VerifyClient(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, servicePublicID, result.ServiceID)
}

func TestVerifyClient_FailsWithWrongSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	clientAppID := vo.NewClientAppPublicID(uuid.New())
	serviceID := vo.NewServiceID(1)

	wrongSecretBytes := []byte("wrong-secret")
	hashedSecret := []byte("hashed-from-correct-secret")

	clientApp := entity.DraftClientApp(serviceID, "Test Client", nil, hashedSecret)

	// Create mocks
	mockClientAppQueryRepo := mock.NewMockClientAppQueryRepository(ctrl)
	mockServiceQueryRepo := mock.NewMockServiceQueryRepository(ctrl)
	mockHasher := mock.NewMockHasher(ctrl)

	usecase := NewUsecase(
		mockClientAppQueryRepo,
		mockServiceQueryRepo,
		mockHasher,
	)

	req := &in.VerifyClientRequest{
		ClientID:     clientAppID,
		ClientSecret: wrongSecretBytes,
	}

	// Set expectations
	mockClientAppQueryRepo.EXPECT().FindByPublicID(ctx, clientAppID).Return(clientApp, nil)

	// Compare should fail with wrong secret
	mockHasher.EXPECT().Compare(ctx, hashedSecret, wrongSecretBytes).Return(assert.AnError)

	// Execute
	result, err := usecase.VerifyClient(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
