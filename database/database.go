package database

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/viniblima/zpe/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

/*
Esta função tem por objetivo iniciar o banco de dados com as devidas tabelas
e inicia alguns preenchimentos necessários. Seguindo os seguintes passos:

1. Conexao com o Banco de Dados;
2. Rodar as migrações e, assim, iniciando a(s) tabela(s);
3. Verificacao se há um usuário com o perfil de administrador para que
o sistema não caia em um loop em que não há nenhum usuário que tenha
pelo menos um perfil de modificacao de outros usuários;
*/
func ConnectDb() {

	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")
	db.AutoMigrate(
		&models.User{},
		&models.Role{},
	)

	var roles []models.Role

	db.Where(models.Role{
		Level: 1,
	}).Or(models.Role{
		Level: 2,
	}).Or(models.Role{
		Level: 3,
	}).Find(&roles)

	roleFound := slices.IndexFunc(roles, func(m models.Role) bool {
		return m.Level == 1
	})

	if roleFound == -1 {
		newRole := models.Role{
			Name:  "Admin",
			Level: 1,
		}

		db.Create(&newRole)

		var usersAdmin []models.User
		db.Model(models.User{}).Where(models.Role{
			Level: 1,
		}).Association("Roles").Find(&usersAdmin)

		if len(usersAdmin) == 0 {
			newUser := models.User{
				Name:     os.Getenv("SUPERUSER_NAME"),
				Email:    os.Getenv("SUPERUSER_EMAIL"),
				Password: os.Getenv("SUPERUSER_PASSWORD"),
			}

			db.Create(&newUser)
			db.Model(&newUser).Omit("Roles.*").Association("Roles").Append(&newRole)
		}
	}

	roleFound = slices.IndexFunc(roles, func(m models.Role) bool {
		return m.Level == 2
	})

	if roleFound == -1 {
		newRole := models.Role{
			Name:  "Modifier",
			Level: 2,
		}
		db.Create(&newRole)
	}

	roleFound = slices.IndexFunc(roles, func(m models.Role) bool {
		return m.Level == 3
	})

	if roleFound == -1 {
		newRole := models.Role{
			Name:  "Watcher",
			Level: 3,
		}
		db.Create(&newRole)
	}

	Db = db
}
