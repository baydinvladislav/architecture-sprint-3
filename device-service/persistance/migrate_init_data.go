package persistance

import (
	"gorm.io/gorm"
	"log"
)

func ApplyMigrations(db *gorm.DB) error {
	log.Printf("Starting migrate init application data")

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("Failed to create uuid-ossp extension: %v", err)
		return err
	}

	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_enum') THEN
				CREATE TYPE status_enum AS ENUM ('INSTALL_PENDING', 'INSTALL_COMPLETED', 'INSTALL_FAILED', 'UNINSTALL');
			END IF;
		END$$;
	`).Error; err != nil {
		log.Fatalf("Failed to create status_enum type: %v", err)
		return err
	}

	log.Printf("Starting database migration for models")
	if err := db.AutoMigrate(
		ModuleModel{},
		HouseModuleModel{},
		HouseModuleHistoryStateModel{},
		DeviceModel{},
	); err != nil {
		log.Fatalf("Failed to migrate models database: %v", err)
		return err
	}
	log.Printf("Database models migration run successfully")

	log.Printf("Adding initial rows to the database")
	if err := addInitialModules(db); err != nil {
		log.Fatalf("Failed to add initial modules: %v", err)
		return err
	}

	log.Printf("Initial modules added successfully")
	return nil
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

	log.Printf("Rows in `modules` table created successfully: %d", len(modules))

	for _, module := range modules {
		var existingModule ModuleModel
		if err := db.Where("type = ?", module.Type).First(&existingModule).Error; err == nil {
			continue
		}

		if err := db.Create(&module).Error; err != nil {
			return err
		}

		var devices []DeviceModel
		switch module.Type {
		case heatingModule:
			devices = []DeviceModel{
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
			devices = []DeviceModel{
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
			devices = []DeviceModel{
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
			devices = []DeviceModel{
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
			devices = []DeviceModel{
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

		log.Printf("Created module '%s' with count devices: %d", module.Type, len(devices))
	}

	return nil
}
