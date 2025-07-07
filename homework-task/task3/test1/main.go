package main

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	dsn := "root:123456@tcp(localhost:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&Student{})
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Transaction{})

	name := "张三"
	grade := "三年级"
	name2 := "张三2"
	grade2 := "三年级2"
	name3 := "张三3"
	grade3 := "三年级3"
	stuList := []Student{{Name: &name, Age: 20, Grade: &grade},
		{Name: &name2, Age: 12, Grade: &grade2},
		{Name: &name3, Age: 12, Grade: &grade3}}

	db.Debug().Create(&stuList)
	studentList := []Student{}
	db.Debug().Where("age > ?", 18).Find(&studentList)
	fmt.Println("Student list :", studentList)

	grade = "四年级"
	db.Debug().Model(&Student{}).Where("name = ?", name).Updates(Student{Grade: &grade})
	db.Debug().Model(&Student{}).Where("age < ?", 15).Delete(&Student{})

	db.Debug().Save(&Account{ID: 1, Balance: 250 * 100})
	db.Debug().Save(&Account{ID: 2, Balance: 0 * 100})

	go func() {
		err := TransferTx(db, 1, 2, 200*100)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := TransferTx(db, 1, 2, 200*100)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(3 * time.Second)
}

type Student struct {
	gorm.Model
	ID    uint64 `gorm:"primaryKey;autoIncrement"`
	Name  *string
	Age   uint8
	Grade *string
}

type Account struct {
	gorm.Model
	ID      uint64 `gorm:"primaryKey;autoIncrement"`
	Balance uint64
}

type Transaction struct {
	gorm.Model
	ID            uint64 `gorm:"primaryKey;autoIncrement"`
	FromAccountID uint64 `gorm:"index;not null"`
	ToAccountID   uint64 `gorm:"index;not null"`
	Amount        uint64 `gorm:"not null"`

	FromAccount Account `gorm:"foreignKey:FromAccountID;references:ID;"`
	ToAccount   Account `gorm:"foreignKey:ToAccountID;references:ID;"`
}

func TransferTx(db *gorm.DB, fromID uint64, toID uint64, amount uint64) error {
	return db.Debug().Transaction(func(tx *gorm.DB) error {

		var fromAccount Account
		if err := tx.Debug().Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("balance").
			First(&fromAccount, fromID).Error; err != nil {
			return err
		}

		if fromAccount.Balance < amount {
			return errors.New("insufficient balance")
		}

		if err := tx.Debug().Model(&Account{}).
			Where("id = ? and balance >= ?", fromID, amount).
			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		if err := tx.Debug().Model(&Account{}).
			Where("id = ?", toID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		txRecord := Transaction{
			FromAccountID: fromID,
			ToAccountID:   toID,
			Amount:        amount,
		}
		if err := tx.Debug().Create(&txRecord).Error; err != nil {
			return err
		}

		return nil
	})
}
