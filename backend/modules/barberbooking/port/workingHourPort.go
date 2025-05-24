package barberBookingPort

import(
	"context"	
	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type IWorkingHourService interface{
	GetWorkingHours(ctx context.Context, branchID uint) ([]barberBookingModels.WorkingHour, error) 
	UpdateWorkingHours(ctx context.Context, branchID uint, input []barberBookingDto.WorkingHourInput) error
	CreateWorkingHours(ctx context.Context, branchID uint, input barberBookingDto.WorkingHourInput) error
}
