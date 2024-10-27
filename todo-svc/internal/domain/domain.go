package domain

type TodoList struct {
	Id          int64      `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Items       []TodoItem `gorm:"foreignKey:ListId"`
}

type UsersList struct {
	Id     int64 `json:"id" gorm:"primaryKey"`
	UserId int64 `json:"user_id"`
	ListId int64 `json:"list_id"`
}

type TodoItem struct {
	Id          int64  `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
	ListId      int64  `json:"list_id"` // зв'язок з TodoList
}

type ListsItem struct {
	Id     int64 `json:"id" gorm:"primaryKey"`
	ListId int64 `json:"list_id"`
	ItemId int64 `json:"item_id"`
}
