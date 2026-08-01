package main

import (
	"ariga.io/atlas-provider-gorm/gormschema"
	"fmt"
	"github.com/usesnipet/snipet/internal/model"
	"io"
	"os"
)

func main() {
	stmts, err := gormschema.New("postgres", gormschema.WithModelPosition(map[any]string{
		&model.APIKey{}:                "/home/ag0020/workspace/mayron/snipet/internal/model/api-key.go:7",
		&model.Agent{}:                 "/home/ag0020/workspace/mayron/snipet/internal/model/agent.go:9",
		&model.AgentToKnowledge{}:      "/home/ag0020/workspace/mayron/snipet/internal/model/agent.go:53",
		&model.AgentToKnowledgeIndex{}: "/home/ag0020/workspace/mayron/snipet/internal/model/agent_to_knowledge_index.go:3",
		&model.AgentToLLM{}:            "/home/ag0020/workspace/mayron/snipet/internal/model/agent.go:40",
		&model.Client{}:                "/home/ag0020/workspace/mayron/snipet/internal/model/client.go:18",
		&model.ClientToUser{}:          "/home/ag0020/workspace/mayron/snipet/internal/model/user.go:19",
		&model.Execution{}:             "/home/ag0020/workspace/mayron/snipet/internal/model/execution.go:10",
		&model.ExecutionMessage{}:      "/home/ag0020/workspace/mayron/snipet/internal/model/execution-message.go:7",
		&model.IndexedKnowledgeItem{}:  "/home/ag0020/workspace/mayron/snipet/internal/model/indexed_knowledge_item.go:19",
		&model.Knowledge{}:             "/home/ag0020/workspace/mayron/snipet/internal/model/knowledge.go:18",
		&model.KnowledgeIndex{}:        "/home/ag0020/workspace/mayron/snipet/internal/model/knowledge-index.go:5",
		&model.KnowledgeItem{}:         "/home/ag0020/workspace/mayron/snipet/internal/model/knowledge-item.go:10",
		&model.LLM{}:                   "/home/ag0020/workspace/mayron/snipet/internal/model/llm.go:7",
		&model.RefreshToken{}:          "/home/ag0020/workspace/mayron/snipet/internal/model/refresh-token.go:9",
		&model.Session{}:               "/home/ag0020/workspace/mayron/snipet/internal/model/session.go:7",
		&model.User{}:                  "/home/ag0020/workspace/mayron/snipet/internal/model/user.go:7",
		&model.UserToSession{}:         "/home/ag0020/workspace/mayron/snipet/internal/model/user.go:28",
	})).Load(
		&model.APIKey{},
		&model.Agent{},
		&model.AgentToKnowledge{},
		&model.AgentToKnowledgeIndex{},
		&model.AgentToLLM{},
		&model.Client{},
		&model.ClientToUser{},
		&model.Execution{},
		&model.ExecutionMessage{},
		&model.IndexedKnowledgeItem{},
		&model.Knowledge{},
		&model.KnowledgeIndex{},
		&model.KnowledgeItem{},
		&model.LLM{},
		&model.RefreshToken{},
		&model.Session{},
		&model.User{},
		&model.UserToSession{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
