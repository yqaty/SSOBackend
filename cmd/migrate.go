package cmd

import (
	"fmt"
	"github.com/UniqueStudio/UniqueSSOBackend/config"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/constants"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/core"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/model"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

// 迁移hr的数据 ---> sso

var (
	migrateCmd = &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			//	runSSOMigrate()
			runSSOMigrateRole()
		},
	}
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}
func getRecruitmentDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s sslmode=disable TimeZone=Asia/Shanghai ",
		"81.70.253.156", "postgres", "recruitment", "5432", "g49pftzdpwdtbba9r3mqb2")
	var err error
	recruitmentDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("connect to db error, %v", err))
	}
	sqlDB, err := recruitmentDB.DB()
	if err != nil {
		panic(fmt.Sprintf("get db error, %v", err))
	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(6000 * time.Second)
	return recruitmentDB
}

type HrMembers struct {
	ID           string    `gorm:"column:id;primaryKey"`
	CreatedAt    time.Time `gorm:"column:createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt"`
	Name         string    `gorm:"column:name"`
	Phone        string    `gorm:"column:phone"`
	Mail         string    `gorm:"column:mail"`
	Gender       int       `gorm:"column:gender"`
	WeChatID     string    `gorm:"column:weChatID"`
	JoinTime     string    `gorm:"column:joinTime"`
	IsCaptain    bool      `gorm:"column:isCaptain;type:boolean"`
	IsAdmin      bool      `gorm:"column:isAdmin;type:boolean"`
	Group        string    `gorm:"column:group"`
	Avatar       string    `gorm:"column:avatar"`
	PasswordSalt string    `gorm:"column:passwordSalt"`
	PasswordHash string    `gorm:"column:passwordHash"`
}

func (*HrMembers) TableName() string {
	return "members"
}
func runSSOMigrateRole() {
	recruitmentDB := getRecruitmentDB()
	err := config.Setup(cfgFile)
	if err != nil {
		panic(err)
	}
	if err := core.Setup(); err != nil {
		log.Println("sso database connection failed", err)
		panic(err)
	}
	ssoDB := core.DB

	var hrMembers []HrMembers
	if err := recruitmentDB.Find(&hrMembers).Error; err != nil {
		log.Println("get hr members failed", err)
		panic(err)
	}
	log.Println("find hr members: ", len(hrMembers))
	// member 权限 member
	// admin 权限 member admin

	var userRole model.UserRole
	if err := ssoDB.AutoMigrate(&userRole); err != nil {
		log.Fatalf("auto migrate user role failed: %v", err)
	}
	for _, hrMember := range hrMembers {
		//if hrMember.IsAdmin {
		//	log.Printf("hr member: %s", hrMember.Name)
		//}
		//if hrMember.IsAdmin {
		//	userRole.CreatedAt = hrMember.CreatedAt
		//	userRole.UpdatedAt = hrMember.UpdatedAt
		//	userRole.RoleName = "admin"
		//	userRole.UID = hrMember.ID
		//	if err := ssoDB.Create(&userRole).Error; err != nil {
		//		log.Fatalf("insert user %#v role failed: %v", userRole, err)
		//	}
		//}
		userRole.CreatedAt = hrMember.CreatedAt
		userRole.UpdatedAt = hrMember.UpdatedAt
		userRole.RoleName = "member"
		userRole.UID = hrMember.ID

		if err := ssoDB.Create(&userRole).Error; err != nil {
			log.Fatalf("insert user %#v role failed: %v", userRole, err)
		}
		//break
	}

	log.Println("insert success")
	// candidate 权限 candidate
	// 待迁移
}

func runSSOMigrate() {
	recruitmentDB := getRecruitmentDB()
	err := config.Setup(cfgFile)
	if err != nil {
		panic(err)
	}
	if err := core.Setup(); err != nil {
		log.Println("sso database connection failed", err)
		panic(err)
	}
	ssoDB := core.DB

	//if err := recruitmentDB.Where("\"isAdmin\" = ?", true).First(&hrMember).Error; err != nil {
	//	log.Println("get hr member failed", err)
	//	panic(err)
	//}
	//log.Println("hr member: ", hrMember)
	// 批量查询hr的数据
	var hrMembers []HrMembers
	if err := recruitmentDB.Find(&hrMembers).Error; err != nil {
		log.Println("get hr members failed", err)
		panic(err)
	}
	var insertUsers []model.User
	if err := ssoDB.Find(&insertUsers).Error; err != nil {
		log.Println("get users failed", err)
		panic(err)
	}
	log.Println("insert users: ", len(insertUsers))
	hashmp := make(map[string]int)
	for _, user := range insertUsers {
		hashmp[user.UID] = 1
	}

	var wg sync.WaitGroup
	for i := 0; i < len(hrMembers); i++ {
		//log.Printf("hr member: %#v", hrMember)
		// 逐个插入到sso的数据库中
		if _, ok := hashmp[hrMembers[i].ID]; ok {
			continue
		}
		if hrMembers[i].Phone == "18856936949" { //tmy
			continue
		}
		var user model.User
		user.UID = hrMembers[i].ID
		user.CreatedAt = hrMembers[i].CreatedAt
		user.UpdatedAt = hrMembers[i].UpdatedAt
		user.Name = hrMembers[i].Name
		user.Email = hrMembers[i].Mail
		if hrMembers[i].Mail == "" {
			user.Email = ""
		}
		user.Phone = hrMembers[i].Phone
		user.JoinTime = hrMembers[i].JoinTime
		user.Groups = []string{hrMembers[i].Group}
		user.AvatarURL = hrMembers[i].Avatar
		user.Gender = constants.Gender(hrMembers[i].Gender)

		wg.Add(1)
		log.Printf("user: %#v", user)
		if err := ssoDB.Create(&user).Error; err != nil {
			log.Println("insert user failed", err)
		}
	}
	wg.Wait()
}

func fixGender(hrGender int) int {
	switch hrGender {
	case 1:
		return 1
	case 2:
		return 2
	default:
		return 3
	}
}
