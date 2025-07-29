package barberBookingPort

import(
	"context"	
	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type IWorkingHourService interface{
	GetWorkingHours(ctx context.Context, branchID uint,tenantID uint) ([]barberBookingModels.WorkingHour, error) 
	GetAvailableSlots(
		ctx context.Context,
		branchID uint,
		tenantID uint,
		filter string,       
		fromTime *string,      
		toTime *string,
	) (map[string][]string, error)
	UpdateWorkingHours(ctx context.Context, branchID uint,tenantID uint ,input []barberBookingDto.WorkingHourInput) error
	CreateWorkingHours(ctx context.Context, branchID uint, input barberBookingDto.WorkingHourInput) error
}
