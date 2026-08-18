package model

type AgentToKnowledgeIndex struct {
	AgentID string `gorm:"type:uuid;primaryKey"`
	IndexID string `gorm:"type:uuid;primaryKey"`

	Agent *Agent          `gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE" json:"agent"`
	Index *KnowledgeIndex `gorm:"foreignKey:IndexID;references:ID;constraint:OnDelete:CASCADE" json:"index"`
}
