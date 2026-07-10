package milvus

import (
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

func SourceSchema(collectionName string, dim int) *entity.Schema {
	schema := entity.NewSchema().WithName(collectionName)
	for _, field := range SourceFields(dim) {
		schema = schema.WithField(field)
	}
	for _, fn := range SourceFunctions() {
		schema = schema.WithFunction(fn)
	}
	return schema
}

func SourceFields(dim int) []*entity.Field {
	return []*entity.Field{
		entity.NewField().
			WithName("id").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(36).
			WithIsPrimaryKey(true),
		entity.NewField().
			WithName("seqId").
			WithDataType(entity.FieldTypeInt64),
		entity.NewField().
			WithName("indexedItemId").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(36),
		entity.NewField().
			WithName("dense").
			WithDataType(entity.FieldTypeFloatVector).
			WithDim(int64(dim)),
		entity.NewField().
			WithName("sparse").
			WithDataType(entity.FieldTypeSparseVector),
		entity.NewField().
			WithName("content").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(8192).
			WithEnableAnalyzer(true),
		entity.NewField().
			WithName("createdAt").
			WithDataType(entity.FieldTypeInt64),
		entity.NewField().
			WithName("updatedAt").
			WithDataType(entity.FieldTypeInt64).
			WithNullable(true),
		entity.NewField().
			WithName("metadata").
			WithDataType(entity.FieldTypeJSON).
			WithNullable(true),
	}
}

func SourceIndexSchema(collectionName string) []milvusclient.CreateIndexOption {
	return []milvusclient.CreateIndexOption{
		milvusclient.NewCreateIndexOption(
			collectionName,
			"dense",
			index.NewIvfFlatIndex(entity.IP, 128),
		).WithIndexName("dense_idx"),
		milvusclient.NewCreateIndexOption(
			collectionName,
			"sparse",
			index.NewGenericIndex("", map[string]string{
				index.IndexTypeKey:    string(index.SparseInverted),
				index.MetricTypeKey:   string(entity.BM25),
				"inverted_index_algo": "DAAT_MAXSCORE",
			}),
		),
	}
}

func SourceFunctions() []*entity.Function {
	bm25Func := entity.NewFunction().
		WithName("bm25_emb").
		WithType(entity.FunctionTypeBM25).
		WithInputFields("content").
		WithOutputFields("sparse")
	bm25Func.Description = "bm25 function"
	return []*entity.Function{bm25Func}
}
