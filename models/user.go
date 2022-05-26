package models

type User struct {
	ID        uint   `json:"id" gorm:"unique_index"`
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique_index"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Password  []byte `json:"-"`
}

type Post struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	AuthorID  uint   `json:"author_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt *int64 `json:"deleted_at"`
}
