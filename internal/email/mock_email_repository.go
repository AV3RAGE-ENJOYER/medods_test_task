package email

import (
	"context"
	"log"

	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/service"
)

type MockEmailRepository struct{}

var _ service.EmailRepository = &MockEmailRepository{}

func (em *MockEmailRepository) NotifyUser(ctx context.Context, email string) error {
	log.Printf(" [Warning] Different ip address for %s", email)
	return nil
}
