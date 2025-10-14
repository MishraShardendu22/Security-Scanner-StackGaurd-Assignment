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
	RFilename   string `json:"rfilename" bson:"rfilename"`
	FileContent string `json:"file_content" bson:"file_content"`
}

type AI_REQUEST struct {
	mgm.DefaultModel `bson:",inline"`
	RequestID        string `json:"request_id" bson:"request_id"`

	Siblings    []SIBLING    `json:"siblings" bson:"siblings"`
	Discussions []DISCUSSION `json:"discussions" bson:"discussions"`
}

type Finding struct {
	SecretType string `json:"secret_type" bson:"secret_type"`
	Pattern    string `json:"pattern" bson:"pattern"`
	Secret     string `json:"secret" bson:"secret"`

	SourceType string `json:"source_type" bson:"source_type"`

	FileName string `json:"file_name,omitempty" bson:"file_name,omitempty"`
	Line     int    `json:"line,omitempty" bson:"line,omitempty"`

	DiscussionNum   int64  `json:"discussion_num,omitempty" bson:"discussion_num,omitempty"`
	DiscussionTitle string `json:"discussion_title,omitempty" bson:"discussion_title,omitempty"`
	DiscussionRepo  string `json:"discussion_repo,omitempty" bson:"discussion_repo,omitempty"`
}

type SCANNED_RESOURCE struct {
	Type     string    `json:"type" bson:"type"`
	ID       string    `json:"id" bson:"id"`
	Findings []Finding `json:"findings" bson:"findings"`
}

type SCAN_RESULT struct {
	mgm.DefaultModel `bson:",inline"`
	RequestID        string             `json:"request_id" bson:"request_id"`
	ScannedResources []SCANNED_RESOURCE `json:"scanned_resources" bson:"scanned_resources"`
}
