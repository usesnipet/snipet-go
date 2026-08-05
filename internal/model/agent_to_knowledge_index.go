package model

type AgentToKnowledgeIndex struct {
	AgentID string `gorm:"type:uuid;primaryKey"`
	IndexID string `gorm:"type:uuid;primaryKey"`

	Agent Agent          `gorm:"foreignKey:AgentID"`
	Index KnowledgeIndex `gorm:"foreignKey:IndexID"`
}
