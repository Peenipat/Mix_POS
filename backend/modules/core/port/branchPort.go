package corePort

type CreateBranchInput struct {
	Name    string `json:"name"    example:"My New Store"`
	Address string `json:"address" example:"123 Main St."`
  }

type UpdateBranchInput struct {
    Name string `json:"name" example:"สำนักงานใหญ่ใหม่"`
}