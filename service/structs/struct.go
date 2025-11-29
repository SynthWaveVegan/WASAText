package structs

type Identifier struct {
	Id string //'json:"Identifier"'
}

type User struct {
	Username string
	UserId   Identifier
}

type Message struct {
	MessageBody string
	Comments    []Comment
	UploaderId  Identifier
	MessageId   Identifier
	Date        string
	MediaType   string
}

type Comment struct {
	CommentId   Identifier
	MessageId   Identifier
	CommentBody string
	Date        string
	UploaderId  Identifier
}

type Group struct {
	GroupId    Identifier
	GroupName  string
	Users      []User
	Messages   []Message
	GroupPhoto Photo
}

type Conversation struct {
	ConversationId Identifier
	UserConnected  Identifier
	UserHosting    Identifier
	ChatName       string
}
