package vo

import "github.com/1340691923/eve-plugin-sdk-go/ev_api/bson"

type MongoExecRes struct {
	Result bson.M `json:"result"`
}

type MongoUpdateRes struct {
	MatchedCount  int64       // The number of documents matched by the filter.
	ModifiedCount int64       // The number of documents modified by the operation.
	UpsertedCount int64       // The number of documents upserted by the operation.
	UpsertedID    interface{} // The _id field of the upserted document, or nil if no upsert was done.
}
