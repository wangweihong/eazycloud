package ielasticsearch

/*
  仅用来解释elasticsearch api的作用
*/
import (
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/dynamicmapping"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/indexcheckonstartup"
)

// From github.com/elastic/go-elasticsearch/v8/typedapi/types/indexsettings.go
type IndexSettings struct {
	// 定义索引的分析器
	Analysis *types.IndexSettingsAnalysis `json:"analysis,omitempty"`
	// Analyze Settings to define analyzers, tokenizers, token filters and character
	// filters.
	Analyze            *types.SettingsAnalyze                   `json:"analyze,omitempty"`
	AutoExpandReplicas *string                                  `json:"auto_expand_replicas,omitempty"`
	Blocks             *types.IndexSettingBlocks                `json:"blocks,omitempty"`
	CheckOnStartup     *indexcheckonstartup.IndexCheckOnStartup `json:"check_on_startup,omitempty"`
	Codec              *string                                  `json:"codec,omitempty"`
	CreationDate       types.StringifiedEpochTimeUnitMillis     `json:"creation_date,omitempty"`
	CreationDateString types.DateTime                           `json:"creation_date_string,omitempty"`
	DefaultPipeline    *string                                  `json:"default_pipeline,omitempty"`
	FinalPipeline      *string                                  `json:"final_pipeline,omitempty"`
	Format             string                                   `json:"format,omitempty"`
	GcDeletes          types.Duration                           `json:"gc_deletes,omitempty"`
	Hidden             string                                   `json:"hidden,omitempty"`
	Highlight          *types.SettingsHighlight                 `json:"highlight,omitempty"`
	Index              *IndexSettings                           `json:"index,omitempty"`
	IndexSettings      map[string]json.RawMessage               `json:"-"`
	// IndexingPressure Configure indexing back pressure limits.
	IndexingPressure              *types.IndicesIndexingPressure `json:"indexing_pressure,omitempty"`
	IndexingSlowlog               *types.IndexingSlowlogSettings `json:"indexing.slowlog,omitempty"`
	Lifecycle                     *types.IndexSettingsLifecycle  `json:"lifecycle,omitempty"`
	LoadFixedBitsetFiltersEagerly *bool                          `json:"load_fixed_bitset_filters_eagerly,omitempty"`
	// Mapping Enable or disable dynamic mapping for an index.
	Mapping                 *types.MappingLimitSettings `json:"mapping,omitempty"`
	MaxDocvalueFieldsSearch *int                        `json:"max_docvalue_fields_search,omitempty"`
	MaxInnerResultWindow    *int                        `json:"max_inner_result_window,omitempty"`
	MaxNgramDiff            *int                        `json:"max_ngram_diff,omitempty"`
	MaxRefreshListeners     *int                        `json:"max_refresh_listeners,omitempty"`
	MaxRegexLength          *int                        `json:"max_regex_length,omitempty"`
	MaxRescoreWindow        *int                        `json:"max_rescore_window,omitempty"`
	MaxResultWindow         *int                        `json:"max_result_window,omitempty"`
	MaxScriptFields         *int                        `json:"max_script_fields,omitempty"`
	MaxShingleDiff          *int                        `json:"max_shingle_diff,omitempty"`
	MaxSlicesPerScroll      *int                        `json:"max_slices_per_scroll,omitempty"`
	MaxTermsCount           *int                        `json:"max_terms_count,omitempty"`
	Merge                   *types.Merge                `json:"merge,omitempty"`
	Mode                    *string                     `json:"mode,omitempty"`
	// 指定了每个主分片的副本数量。
	//	- 副本是主分片的复制，用于提供数据冗余和容错性
	NumberOfReplicas      string `json:"number_of_replicas,omitempty"`
	NumberOfRoutingShards *int   `json:"number_of_routing_shards,omitempty"`
	// 指定了索引的主分片数量。
	//	- 主分片负责存储索引的数据，
	//  - 也负责处理搜索请求
	NumberOfShards       string                     `json:"number_of_shards,omitempty"`
	Priority             string                     `json:"priority,omitempty"`
	ProvidedName         *string                    `json:"provided_name,omitempty"`
	Queries              *types.Queries             `json:"queries,omitempty"`
	QueryString          *types.SettingsQueryString `json:"query_string,omitempty"`
	RefreshInterval      types.Duration             `json:"refresh_interval,omitempty"`
	Routing              *types.IndexRouting        `json:"routing,omitempty"`
	RoutingPartitionSize types.Stringifiedinteger   `json:"routing_partition_size,omitempty"`
	RoutingPath          []string                   `json:"routing_path,omitempty"`
	Search               *types.SettingsSearch      `json:"search,omitempty"`
	Settings             *IndexSettings             `json:"settings,omitempty"`
	// Similarity Configure custom similarity settings to customize how search results are
	// scored.
	Similarity  map[string]types.SettingsSimilarity `json:"similarity,omitempty"`
	SoftDeletes *types.SoftDeletes                  `json:"soft_deletes,omitempty"`
	Sort        *types.IndexSegmentSort             `json:"sort,omitempty"`
	// Store The store module allows you to control how index data is stored and accessed
	// on disk.
	Store               *types.Storage                 `json:"store,omitempty"`
	TimeSeries          *types.IndexSettingsTimeSeries `json:"time_series,omitempty"`
	TopMetricsMaxSize   *int                           `json:"top_metrics_max_size,omitempty"`
	Translog            *types.Translog                `json:"translog,omitempty"`
	Uuid                *string                        `json:"uuid,omitempty"`
	VerifiedBeforeClose string                         `json:"verified_before_close,omitempty"`
	Version             *types.IndexVersioning         `json:"version,omitempty"`
}

type TypeMapping struct {
	AllField             *types.AllField                    `json:"all_field,omitempty"`
	DataStreamTimestamp_ *types.DataStreamTimestamp         `json:"_data_stream_timestamp,omitempty"`
	DateDetection        *bool                              `json:"date_detection,omitempty"`
	Dynamic              *dynamicmapping.DynamicMapping     `json:"dynamic,omitempty"`
	DynamicDateFormats   []string                           `json:"dynamic_date_formats,omitempty"`
	DynamicTemplates     []map[string]types.DynamicTemplate `json:"dynamic_templates,omitempty"`
	Enabled              *bool                              `json:"enabled,omitempty"`
	FieldNames_          *types.FieldNamesField             `json:"_field_names,omitempty"`
	IndexField           *types.IndexField                  `json:"index_field,omitempty"`
	Meta_                types.Metadata                     `json:"_meta,omitempty"`
	NumericDetection     *bool                              `json:"numeric_detection,omitempty"`
	// 指定索引中的字段映射
	//"properties": {
	//			"productName": {
	//				"type": "text",
	//				"analyzer": "standard"
	//			},
	//			"annual_rate": {
	//				"type": "keyword"
	//			},
	//			"describe": {
	//				"type": "text",
	//				"analyzer": "standard"
	//			}
	//		}
	// productName:指定了一个名为 "productName" 的字段，其类型为 "text"。
	//		- "text"类型适用于全文搜索，会被分析器处理以便建立倒排索引。
	//		- "analyzer": "standard",指定分析器为"standard"分析器，该分析器会按照标准的分词规则进行分析。
	// annual_rate：指定了一个名为 "annual_rate" 的字段，其类型为 "keyword"。
	//		- "keyword"类型适用于精确匹配和聚合操作，不会被分析器处理。
	// describe：指定了一个名为"describe"的字段，
	//		- 其类型为 "text"，使用了与 "productName" 字段相同的 "standard" 分析器。
	Properties map[string]types.Property     `json:"properties,omitempty"`
	Routing_   *types.RoutingField           `json:"_routing,omitempty"`
	Runtime    map[string]types.RuntimeField `json:"runtime,omitempty"`
	Size_      *types.SizeField              `json:"_size,omitempty"`
	Source_    *types.SourceField            `json:"_source,omitempty"`
	Subobjects *bool                         `json:"subobjects,omitempty"`
}

type SearchResponse struct {
	// 搜索请求花费的时间,以毫秒为单位。
	Took int `json:"took"`
	// 指示搜索是否超时
	TimedOut bool `json:"timed_out"`
	// 关于搜索请求在Elasticsearch集群中各分片（shard）的执行情况的统计信息。
	Shards struct {
		// 总共涉及的分片数
		Total int `json:"total"`
		// 成功执行搜索请求的分片数。
		Successful int `json:"successful"`
		// 跳过的分片数
		Skipped int `json:"skipped"`
		// 执行搜索请求失败的分片数
		Failed int `json:"failed"`
	} `json:"_shards"`
	// 搜索结果的详细信息
	Hits struct {
		// 搜索结果的总数
		Total struct {
			// 搜索结果的数量
			Value int `json:"value"`
			// 与搜索结果数量的关系，可以是 eq（等于），gte（大于等于）等
			Relation string `json:"relation"`
		} `json:"total"`
		// 搜索结果中最高的分数
		MaxScore float64 `json:"max_score"`
		// 搜索命中的文档列表
		Hits []struct {
			// 文档所在的索引名称
			Index string `json:"_index"`
			// 文档的ID
			ID string `json:"_id"`
			// 文档的匹配分数。分数越高，表示文档越匹配搜索条件
			Score float64 `json:"_score"`
			// 文档的源数据，包含了文档的实际内容
			Source map[string]any `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type HitData struct {
	// 文档所在的索引名称
	Index string `json:"_index"`
	// 文档的ID
	ID string `json:"_id"`
	// 文档的匹配分数。分数越高，表示文档越匹配搜索条件
	Score float64 `json:"_score"`
	// 文档的源数据，包含了文档的实际内容
	// 通过`"_source": false`可以不返回此字段
	Source map[string]any `json:"_source"`
	// 特定字段的返回值
	// 搜索条件包含`{"fields": ["id","name"]}`时，返回id/name字段的数据
	Fields map[string][]any `json:"fields"`
}

/*
{
	"settings": {
		"number_of_shards": 5,
		"number_of_replicas": 1,
		"analysis": {
			"analyzer": {
				"es_std": {
					"type": "standard",
					"stopwords": "_english_"
				}
			}
		}
	},
	"mappings": {
		"properties": {
			"productName": {
				"type": "text",
				"analyzer": "standard"
			},
			"annual_rate": {
				"type": "keyword"
			},
			"describe": {
				"type": "text",
				"analyzer": "standard"
			}
		}
	}
}
*/
// 见github.com\elastic\go-elasticsearch\v8\typedapi\indices\create\create.go
type IndexCreateParam struct {
	// 指定索引的设置
	Settings types.IndexSettings `json:"settings"`
	Mappings types.TypeMapping   `json:"mappings"`
}
