package structs

type Identifier struct {
	Id string   `json:"Identifier"`
}

type User struct {
	Username string       `json:"Name"`
	UserId   Identifier   `json:"userId"`
	UserPhoto string      `json:"UserPhoto"`
}

type Comment struct {
	CommentId   Identifier  `json:"commentId"`
	MessageId   Identifier  `json:"messageId"`
	CommentBody string      `json:"CommentBody"`
	Date        string      `json:"Date"`
	UploaderId  Identifier  `json:"UploaderId"`
}

type Message struct {
	MessageBody string       `json:"MessageBody"`
	Comments    []Comment    `json:"Comments"`
	UploaderId  Identifier   `json:"UploaderId"`
	MessageId   Identifier   `json:"messageId"`
	Date        string       `json:"Date"`
	MediaType   string       `json:"MediaType"`
}

type Group struct {
	GroupId    Identifier  `json:"groupId"`
	GroupName  string      `json:"Name"`
	Users      []User      `json:"Users"`
	Messages   []Message   `json:"Messages"`
}

type Conversation struct {
	ConversationId Identifier  `json:"Identifier"`
	UserConnected  Identifier  `json:"UserConnected"`
	UserHosting    Identifier  `json:"userHosting"`
	ChatName       string      `json:"Name"`
}

