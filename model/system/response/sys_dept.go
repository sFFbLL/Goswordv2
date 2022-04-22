package response

import "project/model/system"

type SysDeptResponse struct {
	Dept system.SysDept `json:"Dept"`
}
