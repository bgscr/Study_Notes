package main

import (
	"fmt"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:123456@tcp(localhost:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	users := []User{
		{Username: "张三", Password: "zhangsan123", Email: "zhangsan@example.com"},
		{Username: "李四", Password: "lisi456", Email: "lisi@example.com"},
		{Username: "王五", Password: "wangwu789", Email: "wangwu@example.com"},
		{Username: "赵六", Password: "zhaoliu123", Email: "zhaoliu@example.com"},
	}
	db.CreateInBatches(users, 2) //

	var posts []Post
	for _, user := range users {
		for i := 1; i <= 3; i++ {
			posts = append(posts, Post{
				Title:   "文章标题-" + user.Username,
				Content: "这是" + user.Username + "的第" + strconv.Itoa(i) + "篇文章",
				UserID:  user.ID,
			})
		}
	}
	db.CreateInBatches(posts, 5)

	var comments []Comment
	for _, post := range posts {
		for i := 0; i < 5; i++ {
			commentUser := users[i%len(users)]
			comments = append(comments, Comment{
				Content: "用户" + commentUser.Username + "的评论",
				UserID:  commentUser.ID,
				PostID:  post.ID,
			})
		}
	}
	db.CreateInBatches(comments, 10)

	var user User
	db.Debug().Preload("Posts.Comments").Preload("Posts.Comments.User").
		Where("id = ?", 4).
		First(&user)

	fmt.Println("find user:", user)

	var post Post
	db.Debug().
		Preload("Comments").Preload("Comments.User").
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count ASC").
		First(&post)
	fmt.Println("find post:", post)

	var deleteComments []Comment
	db.Where("post_id = ?", posts[len(posts)-1].ID).Find(&deleteComments)
	for _, c := range deleteComments {
		db.Delete(&c)
	}
}

type User struct {
	gorm.Model
	Username  string    `gorm:"size:255;not null;unique"`
	Password  string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;not null;unique"`
	PostCount uint64    `gorm:"default:0"`
	Posts     []Post    `gorm:"foreignKey:UserID"`
	Comments  []Comment `gorm:"foreignKey:UserID"`
}

type Post struct {
	gorm.Model
	Title         string `gorm:"size:255;not null"`
	Content       string `gorm:"type:text;not null"`
	UserID        uint
	CommentStatus string    `gorm:"default:'no_comments'"`
	Comments      []Comment `gorm:"foreignKey:PostID"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	UserID  uint
	PostID  uint
	User    User `gorm:"foreignKey:UserID"`
	Post    Post `gorm:"foreignKey:PostID"`
}

func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(&User{}).
		Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + 1"))
	return
}

func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_status", "has_comments").
		Error
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {

	var remainingComments int64
	if err = tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&remainingComments).
		Error; err != nil {
		return
	}

	if remainingComments > 0 {
		return
	}

	return tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_status", "no_comments").
		Error
}
