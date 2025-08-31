package service

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/uchupx/saceri-chatbot-api/internal/database"
	"github.com/uchupx/saceri-chatbot-api/internal/models"
	"github.com/uchupx/saceri-chatbot-api/internal/repository"
	"github.com/uchupx/saceri-chatbot-api/pkg/apierror"
)

type ChatbotService struct {
	SettingRepo repository.SettingRepoInterface
	EventRepo   repository.EventRepoInterface
	Cache       *database.Cache
}

func (s *ChatbotService) UpdatePromptContext(ctx context.Context) *apierror.APIerror {
	prompt, err := s.SettingRepo.GetByKey(ctx, models.SettingKeyPrompts)
	if err != nil {
		return apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	if prompt == nil {
		return apierror.NewAPIError(echo.ErrNotFound.Code, fmt.Errorf("Prompts setting not found"))
	}

	now := time.Now()
	endOfYear := time.Date(now.Year(), 12, 31, 0, 0, 0, 0, now.Location())

	dates := []time.Time{
		now,
		endOfYear,
	}

	events, err := s.EventRepo.GetEventsByDates(ctx, dates)
	if err != nil {
		return apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	eventString := ""
	for idx, event := range events {
		var donationTarget uint
		if event.DonationTarget != nil {
			donationTarget = *event.DonationTarget
		}

		eventString += fmt.Sprintf("%d. %s. \n Tempat: %s. \n Tanggal: %s. \n Deskripsi: %s .\n Target Donasi: Rp.%d \n\n ", idx+1, event.Name, event.Place, event.StartEvent.Format("2006-01-02"), event.Description, donationTarget)
	}

	context := `
%s

Event yang akan datang:
%s
Fokuskan jawaban Anda pada topik ini jika tidak sesuai topik makan hanya menjawab Maaf, pertanyaan anda tidak sesuai topik
	.`

	context = fmt.Sprintf(context, prompt.Value, eventString)

	err = s.Cache.Put(ctx, "prompts_context", context, time.Minute*5)
	if err != nil {
		return apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	return nil
}

func (s *ChatbotService) UpdateStaticToken(ctx context.Context) *apierror.APIerror {
	return nil
}
