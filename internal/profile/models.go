package profile

import "time"

type Profile struct {
	ID         string    `json:"id"`          // UUID профиля (primary key)
	UserID     string    `json:"user_id"`     // id пользователя (foreign key)
	Name       string    `json:"name"`        // "Алексей"
	NickName   string    `json:"nick_name"`   // прозвище, должно быть уникальным
	Gender     string    `json:"gender"`      // "male", "female"
	AgeGroup   string    `json:"age_group"`   // "18-20", "21-25"
	City       string    `json:"city"`        // "Москва"
	Profession string    `json:"profession"`  // "IT", "doctor"
	Smoking    string    `json:"smoking"`     // "none", "sometimes"
	Goal       string    `json:"goal"`        // "dating", "friendship"
	Hobbies    []string  `json:"hobbies"`     // ["travel", "music"]
	SocialLink string    `json:"social_link"` // "tg://username"
	Rating     int64     `json:"rating"`      // при регистрации устанавливается в 10 единиц
	CreatedAt  time.Time `json:"created_at"`  // дата создания
	UpdatedAt  time.Time `json:"updated_at"`  // дата обновления
}
