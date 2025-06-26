package barberBookingDto
import(
	"time"
)

type WorkingHourInput struct {
	Weekday   int       `json:"weekday"`    // 0–6
	StartTime time.Time `json:"start_time"` // เวลาเปิด
	EndTime   time.Time `json:"end_time"`   // เวลาปิด
	IsClosed  bool		`json:"is_closed"` 
}
