package main

import "context"

func (r *GormRepository[M, E]) Insert(ctx context.Context, entity *E) error {
	// отображаем данные из Entity в DTO
	var start M
	model := start.FromEntity(*entity).(M)

	// создаём новую запись в базе данных
	err := r.db.WithContext(ctx).Create(&model).Error
	// обработка ошибки

	// отображаем новую запись из базы в Entity
	*entity = model.ToEntity()
	return nil
}

func (r *GormRepository[M, E]) FindByID(ctx context.Context, id uint) (E, error) {
	// извлекаем запись по id из базы данных
	var model M
	err := r.db.WithContext(ctx).First(&model, id).Error
	// обработка ошибки

	// отображаем запись в Entity
	return model.ToEntity(), nil
}

func (r *GormRepository[M, E]) Find(ctx context.Context, specification Specification) ([]E, error) {
	// получаем записи по некоторому критерию
	var models []M
	err := r.db.WithContext(ctx).Where(specification.GetQuery(), specification.GetValues()...).Find(&models).Error
	// обработка ошибки

	// отображаем все записи в Entities
	result := make([]E, 0, len(models))
	for _, row := range models {
		result = append(result, row.ToEntity())
	}

	return result, nil
}
