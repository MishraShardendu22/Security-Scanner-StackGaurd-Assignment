package models

import "github.com/kamva/mgm/v3"

type MODEL_URL struct {
	mgm.DefaultModel `bson:",inline"`
	URL              string `json:"url" bson:"url"`
}

type BaseAI struct {
	mgm.DefaultModel  `bson:",inline"`
	Org               string `json:"org" bson:"org"`
	User              string `json:"user" bson:"user"`
	Model             string `json:"model" bson:"model"`
	IncludePRS        bool   `json:"include_prs" bson:"include_prs"`
	IncludeDiscussion bool   `json:"include_discussion" bson:"include_discussion"`
}

type AI_Models struct {
	BaseAI
	Model_ID string `json:"model_id" bson:"model_id"`
}

type AI_SPACES struct {
	BaseAI
	Space_ID string `json:"space_id" bson:"space_id"`
}

type AI_DATASETS struct {
	BaseAI
	Dataset_ID string `json:"dataset_id" bson:"dataset_id"`
}

type DISCUSSION struct {
	Num           int64  `json:"num" bson:"num"`
	Title         string `json:"title" bson:"title"`
	Status        string `json:"status" bson:"status"`
	IsPullRequest bool   `json:"isPullRequest" bson:"isPullRequest"`
	CreatedAt     string `json:"createdAt" bson:"createdAt"`
	AuthorName    string `json:"author_name" bson:"author_name"`
	RepoName      string `json:"repo_name" bson:"repo_name"`
	NumComments   int64  `json:"numComments" bson:"numComments"`
	Pinned        bool   `json:"pinned" bson:"pinned"`
}

type SIBLING struct {
	RFilename string `json:"rfilename" bson:"rfilename"`
}

type AI_REQUEST struct {
	mgm.DefaultModel `bson:",inline"`
	RequestID   string       `json:"request_id" bson:"request_id"`
	
	Siblings    []SIBLING    `json:"siblings" bson:"siblings"`
	Discussions []DISCUSSION `json:"discussions" bson:"discussions"`
}
