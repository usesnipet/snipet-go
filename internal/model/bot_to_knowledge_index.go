package model

type BotToKnowledgeIndex struct {
	BotID   string `gorm:"type:uuid;primaryKey"`
	IndexID string `gorm:"type:uuid;primaryKey"`

	Bot   Bot            `gorm:"foreignKey:BotID"`
	Index KnowledgeIndex `gorm:"foreignKey:IndexID"`
}
