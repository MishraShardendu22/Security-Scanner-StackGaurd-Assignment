package models

type MODEL_URL struct {
	URL string `json:"url" bson:"url"`
}

type AI_Models struct {
	Org               string `json:"org" bson:"org"`
	User              string `json:"user" bson:"user"`
	Model_ID          string `json:"model_id" bson:"model_id"`
	Space_ID          string `json:"space_id" bson:"space_id"`
	Dataset_ID        string `json:"dataset_id" bson:"dataset_id"`
	IncludePRS        bool   `json:"include_prs" bson:"include_prs"`
	IncludeDiscussion bool   `json:"include_discussion" bson:"include_discussion"`
}
