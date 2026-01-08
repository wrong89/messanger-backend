package chat

import "time"

type Chat struct {
	ID        uint64
	Type      string
	Address   string
	CreatedAt time.Time
	// ChatMembers
}

// type ChatMembers struct {
// 	Role          string
// 	ChatID        uint64
// 	UserID        []uint64
// 	JoinedAt      time.Time
// 	IsBanned      bool
// 	LastReadMsgID *uint64
// }

// func NewChatMember(role string, chatID uint64, userID []uint64) ChatMembers {
// 	return ChatMembers{
// 		Role:          role,
// 		ChatID:        chatID,
// 		UserID:        userID,
// 		JoinedAt:      time.Now(),
// 		IsBanned:      false,
// 		LastReadMsgID: nil,
// 	}
// }
