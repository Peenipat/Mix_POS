package barberBookingDto

type BranchWorkloadSummary struct {
    TenantID     uint `json:"tenant_id"`
    BranchID     uint `json:"branch_id"`
    NumWorked    int  `json:"num_worked"`     // ช่างที่มาทำงานจริง
    TotalBarbers int  `json:"total_barbers"`  // ช่างทั้งหมดที่ขึ้นทะเบียน
}
