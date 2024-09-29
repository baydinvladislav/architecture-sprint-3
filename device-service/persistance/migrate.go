package persistance

import (
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("failed to create uuid-ossp extension: %v", err)
	}

	if err := db.AutoMigrate(
		ModuleModel{},
		HouseModuleModel{},
		Device{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if err := addInitialModules(db); err != nil {
		log.Fatalf("failed to add initial modules: %v", err)
	}
}

func addInitialModules(db *gorm.DB) error {
	heatingModule := "Модуль обогрева"
	movingModule := "Модуль датчика движения"
	gatesModule := "Модуль автоматических ворот"
	watchingModule := "Модуль системы наблюдения"
	lightingModule := "Модуль управления освещением"

	modules := []ModuleModel{
		{
			Type:        heatingModule,
			Description: "Модуль обогрева для поддержания комфортной температуры в доме.",
		},
		{
			Type:        movingModule,
			Description: "Датчик движения для обеспечения безопасности и мониторинга.",
		},
		{
			Type:        gatesModule,
			Description: "Автоматические ворота для удобного доступа и безопасности.",
		},
		{
			Type:        watchingModule,
			Description: "Система видеонаблюдения для защиты вашего дома.",
		},
		{
			Type:        lightingModule,
			Description: "Управление освещением для экономии энергии и создания атмосферы.",
		},
	}

	for _, module := range modules {
		var existingModule ModuleModel
		if err := db.Where("type = ?", module.Type).First(&existingModule).Error; err == nil {
			continue
		}

		if err := db.Create(&module).Error; err != nil {
			return err
		}

		var devices []Device
		switch module.Type {
		case heatingModule:
			devices = []Device{
				{
					Name:        "Электронный термостат",
					VendorName:  "ТеплоКом",
					Description: "Устройство для управления температурой.",
				},
				{
					Name:        "Настенный обогреватель",
					VendorName:  "Солнечный свет",
					Description: "Обогреватель для комфортной температуры в комнате.",
				},
			}
		case movingModule:
			devices = []Device{
				{
					Name:        "Инфракрасный датчик движения",
					VendorName:  "Безопасность+",
					Description: "Датчик для обнаружения движения.",
				},
				{
					Name:        "Умная сигнализация",
					VendorName:  "Защита+",
					Description: "Сигнализация с функцией уведомления.",
				},
			}
		case gatesModule:
			devices = []Device{
				{
					Name:        "Электрический привод ворот",
					VendorName:  "АвтоДвери",
					Description: "Привод для автоматизации открывания ворот.",
				},
				{
					Name:        "Пульт дистанционного управления",
					VendorName:  "Удобный доступ",
					Description: "Пульт для управления автоматическими воротами.",
				},
			}
		case watchingModule:
			devices = []Device{
				{
					Name:        "IP-камера",
					VendorName:  "Безопасный дом",
					Description: "Камера для видеонаблюдения.",
				},
				{
					Name:        "Запись видео",
					VendorName:  "Хранитель",
					Description: "Устройство для записи и хранения видео.",
				},
			}
		case lightingModule:
			devices = []Device{
				{
					Name:        "Умная лампочка",
					VendorName:  "Светлый дом",
					Description: "Лампочка с управлением через приложение.",
				},
				{
					Name:        "Светодиодный контроллер",
					VendorName:  "ЭкоСвет",
					Description: "Контроллер для управления освещением.",
				},
			}
		}

		for _, device := range devices {
			device.ModuleID = module.ID
			if err := db.Create(&device).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
