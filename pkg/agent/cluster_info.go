package agent

type Cluster struct{
	Primary string `json:"primary"`
	Secondary []string `json:"secondary"`
}