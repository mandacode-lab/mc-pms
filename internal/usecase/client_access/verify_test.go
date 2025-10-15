package clientaccess

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/test/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVerifyClient_ComparesWithRawSecretBytes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// Setup test data
	serviceID, _ := service.NewID(1)
	servicePublicID, _ := service.NewPublicID(uuid.New().String())
	serviceName, _ := service.NewName("Test Service")
	serviceDesc, _ := service.NewDescription("A test service")
	service := service.NewService(
		serviceID,
		servicePublicID,
		serviceName,
		serviceDesc,
		true, // isActive
		time.Now().UTC(),
		time.Now().UTC(),
	)

	clientID, _ := client.NewID(1)
	clientPublicID, _ := client.NewPublicID(uuid.New().String())
	clientName, _ := client.NewName("Test Client")
	clientDesc, _ := client.NewDescription("A test client application")
	secretBytes := []byte("raw-secret-bytes")
	hashedSecret := []byte("hashed-from-raw-bytes")
	encodedSecret := []byte("encoded-secret") // Simulate base64 encoded secret

	now := time.Now().UTC()

	// desc, _ := vo.NewClientAppDescription("")
	// clientApp, _ := entity.DraftClientApp(serviceID, "Test Client", desc, hashedSecret)
	client := client.NewClient(
		clientID,
		clientPublicID,
		clientName,
		clientDesc,
		hashedSecret,
		serviceID,
		true, // isActive
		now,
		now,
	)

	// Create mocks
	mockClientAppQueryRepo := mock.NewMockClientQueryRepository(ctrl)
	mockServiceQueryRepo := mock.NewMockServiceQueryRepository(ctrl)
	mockHasher := mock.NewMockHasher(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)

	usecase := NewUsecase(
		mockClientAppQueryRepo,
		mockServiceQueryRepo,
		mockHasher,
		mockEncoder,
	)

	req := &in.VerifyClientInput{
		ClientID:     clientPublicID,
		ClientSecret: encodedSecret, // User sends base64 encoded secret
	}

	// Set expectations
	mockClientAppQueryRepo.EXPECT().FindByPublicID(ctx, clientPublicID).Return(client, nil)

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

	// Setup test data
	serviceID, _ := service.NewID(1)

	clientID, _ := client.NewID(1)
	clientPublicID, _ := client.NewPublicID(uuid.New().String())
	clientName, _ := client.NewName("Test Client")
	clientDesc, _ := client.NewDescription("A test client application")
	secretBytes := []byte("raw-secret-bytes")
	hashedSecret := []byte("hashed-from-raw-bytes")
	encodedSecret := []byte("encoded-secret") // Simulate base64 encoded secret

	now := time.Now().UTC()

	client := client.NewClient(
		clientID,
		clientPublicID,
		clientName,
		clientDesc,
		hashedSecret,
		serviceID,
		true, // isActive
		now,
		now,
	)

	// Create mocks
	mockClientAppQueryRepo := mock.NewMockClientQueryRepository(ctrl)
	mockServiceQueryRepo := mock.NewMockServiceQueryRepository(ctrl)
	mockHasher := mock.NewMockHasher(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)

	usecase := NewUsecase(
		mockClientAppQueryRepo,
		mockServiceQueryRepo,
		mockHasher,
		mockEncoder,
	)

	req := &in.VerifyClientInput{
		ClientID:     clientPublicID,
		ClientSecret: encodedSecret, // User sends base64 encoded secret
	}

	// Set expectations
	mockClientAppQueryRepo.EXPECT().FindByPublicID(ctx, clientPublicID).Return(client, nil)

	// CRITICAL: Decode the encoded secret to get raw bytes
	mockEncoder.EXPECT().Decode(encodedSecret).Return(secretBytes, nil)

	// CRITICAL: Compare should be called with hashedSecret and decoded raw secretBytes
	mockHasher.EXPECT().Compare(ctx, hashedSecret, secretBytes).Return(assert.AnError)

	// Execute
	result, err := usecase.VerifyClient(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
