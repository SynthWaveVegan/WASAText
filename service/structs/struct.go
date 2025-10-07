package structs

type Identifier struct {

	Identifier           string //'json:"Identifier"'
}

type User struct {

	Username             string
	UserId               Identifier
}

type Message struct {

	MessageBody          string
	Comments             []Comment
	Uploader             Identifier
	MessageId            Identifier
	Date                 string
	ConversationId       Identifier
}

type Comment struct {

	CommentId            Identifier
	MessageId            Identifier
	CommentBody          string
	Date                 string

}

type Photo struct {

	PhotoId              Identifier
	Path                 string
}

type Group struct {

	GroupId              Identifier
	Users                []User
	Messages             []Message
	GroupName            string
	Photo                Photo
}

type Conversation struct {

	ConversationId       Identifier
	UserConnected        User
	Messages             []Message
	ChatName             string


}
