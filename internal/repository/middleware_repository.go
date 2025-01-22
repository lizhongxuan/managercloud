package repository

import (
	"middleware-platform/internal/model"

	"gorm.io/gorm"
)

type MiddlewareRepository struct {
	db *gorm.DB
}

func NewMiddlewareRepository(db *gorm.DB) *MiddlewareRepository {
	// 自动迁移数据库表
	db.AutoMigrate(&model.Middleware{})
	return &MiddlewareRepository{db: db}
}

// FindAll 查找所有中间件
func (r *MiddlewareRepository) FindAll() ([]model.Middleware, error) {
	var middlewares []model.Middleware
	result := r.db.Find(&middlewares)
	return middlewares, result.Error
}

// FindByType 根据类型查找中间件
func (r *MiddlewareRepository) FindByType(middlewareType string) ([]model.Middleware, error) {
	var middlewares []model.Middleware
	result := r.db.Where("type = ?", middlewareType).Find(&middlewares)
	return middlewares, result.Error
}

// Create 创建中间件
func (r *MiddlewareRepository) Create(middleware *model.Middleware) error {
	return r.db.Create(middleware).Error
}

// Update 更新中间件
func (r *MiddlewareRepository) Update(middleware *model.Middleware) error {
	return r.db.Save(middleware).Error
}

// Delete 删除中间件
func (r *MiddlewareRepository) Delete(id uint) error {
	return r.db.Delete(&model.Middleware{}, id).Error
}

func (r *MiddlewareRepository) FindByID(id uint) (*model.Middleware, error) {
	var middleware model.Middleware
	if err := r.db.First(&middleware, id).Error; err != nil {
		return nil, err
	}
	return &middleware, nil
} 