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
	// 3. During verify: user sends encodedSecret in Basic Auth password
	// 4. Usecase: Decode(encodedSecret) -> secretBytes
	// 5. Compare(hashedSecret, secretBytes) should succeed
	secretBytes := []byte("raw-secret-bytes")
	encodedSecret := []byte("base64-encoded-secret")
	hashedSecret := []byte("hashed-from-raw-bytes")

	clientApp := entity.DraftClientApp(serviceID, "Test Client", nil, hashedSecret)

	serviceName, _ := vo.NewServiceName("Test Service")
	service := entity.DraftService(serviceName, nil)
	servicePublicID := service.PublicID()

	// Create mocks
	mockClientAppQueryRepo := mock.NewMockClientAppQueryRepository(ctrl)
	mockServiceQueryRepo := mock.NewMockServiceQueryRepository(ctrl)
	mockHasher := mock.NewMockHasher(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)

	usecase := NewUsecase(
		mockClientAppQueryRepo,
		mockServiceQueryRepo,
		mockHasher,
		mockEncoder,
	)

	req := &in.VerifyClientRequest{
		ClientID:     clientAppID,
		ClientSecret: encodedSecret, // User sends base64 encoded secret
	}

	// Set expectations
	mockClientAppQueryRepo.EXPECT().FindByPublicID(ctx, clientAppID).Return(clientApp, nil)

	// CRITICAL: Decode the encoded secret to get raw bytes
	mockEncoder.EXPECT().Decode(encodedSecret).Return(secretBytes, nil)

	// CRITICAL: Compare should be called with hashedSecret and decoded raw secretBytes
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

	wrongEncodedSecret := []byte("wrong-encoded-secret")
	wrongSecretBytes := []byte("wrong-secret-bytes")
	hashedSecret := []byte("hashed-from-correct-secret")

	clientApp := entity.DraftClientApp(serviceID, "Test Client", nil, hashedSecret)

	// Create mocks
	mockClientAppQueryRepo := mock.NewMockClientAppQueryRepository(ctrl)
	mockServiceQueryRepo := mock.NewMockServiceQueryRepository(ctrl)
	mockHasher := mock.NewMockHasher(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)

	usecase := NewUsecase(
		mockClientAppQueryRepo,
		mockServiceQueryRepo,
		mockHasher,
		mockEncoder,
	)

	req := &in.VerifyClientRequest{
		ClientID:     clientAppID,
		ClientSecret: wrongEncodedSecret,
	}

	// Set expectations
	mockClientAppQueryRepo.EXPECT().FindByPublicID(ctx, clientAppID).Return(clientApp, nil)

	// Decode the wrong encoded secret
	mockEncoder.EXPECT().Decode(wrongEncodedSecret).Return(wrongSecretBytes, nil)

	// Compare should fail with wrong secret
	mockHasher.EXPECT().Compare(ctx, hashedSecret, wrongSecretBytes).Return(assert.AnError)

	// Execute
	result, err := usecase.VerifyClient(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
