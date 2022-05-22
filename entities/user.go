package entities

import (
	"finance/car-finance/back-end/helpers"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `gorm:"column:first_name" json:"first_name,omitempty"`
	LastName  string `gorm:"column:last_name" json:"last_name,omitempty"`
	Code      uint   `gorm:"column:code" json:"code,omitempty"`
	Keyword   string `gorm:"-" json:"keyword"`
}

type RequestUser struct {
	ID        uint   `gorm:"column:id" json:"id,omitempty"`
	FirstName string `gorm:"column:first_name" json:"first_name,omitempty"`
	LastName  string `gorm:"column:last_name" json:"last_name,omitempty"`
	Code      uint   `gorm:"column:code" json:"code,omitempty"`
	Keyword   string `gorm:"-" json:"keyword"`

	PerPage  uint32 `name=per_page,json=perPage" json:"per_page,omitempty"`
	Page     uint32 `name=page" json:"page,omitempty"`
	SortBy   string `name=sort_by,json=sortBy" json:"sort_by,omitempty"`
	SortType string `name=sort_type,json=sortType" json:"sort_type,omitempty"`
}

type UserResponse struct {
	Message string        `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Data    []*UserDetail `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	Meta    *helpers.Meta `protobuf:"bytes,3,opt,name=meta,proto3" json:"meta,omitempty"`
}

type UserDetail struct {
	ID        uint                   `gorm:"column:id" json:"id,omitempty"`
	FirstName string                 `gorm:"column:first_name" json:"first_name,omitempty"`
	LastName  string                 `gorm:"column:last_name" json:"last_name,omitempty"`
	Code      uint                   `gorm:"column:code" json:"code,omitempty"`
	CreatedAt *timestamppb.Timestamp `name=created_at,json=createdAt" json:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `name=updated_at,json=updatedAt" json:"updated_at,omitempty"`
	DeletedAt *timestamppb.Timestamp `name=deleted_at,json=deletedAt" json:"deleted_at,omitempty"`
}

func (t *User) TableName() string {
	return "user"
}

// ToEntity wrap protobuf detail to entity node
func (e *User) ToEntity(req *User) error {
	if req.ID != 0 {
		e.ID = req.ID
	}
	if req.FirstName != "" {
		e.FirstName = req.FirstName
	}
	if req.LastName != "" {
		e.LastName = req.LastName
	}
	if req.Code != 0 {
		e.Code = req.Code
	}
	if req.Keyword != "" {
		e.Keyword = req.Keyword
	}
	return nil
}

// ToProtobuf wrap entitity data to protobuf struct
func (e *User) ToProtobuf() *UserDetail {
	createdAt, _ := ptypes.TimestampProto(e.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(e.UpdatedAt)
	pb := UserDetail{
		ID:        e.ID,
		FirstName: e.FirstName,
		LastName:  e.LastName,
		Code:      e.Code,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	// Id:        uint32(e.ID),
	// 	BrandName: e.BrandName,
	// 	BrandRank: e.BrandRank,
	// 	BrandIcon: e.BrandIcon,
	// 	CreatedAt: createdAt,
	// 	UpdatedAt: updatedAt,

	return &pb
}
