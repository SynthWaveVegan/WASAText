package structs

type Identifier struct {

	Id                   string //'json:"Identifier"'
}

type User struct {

	Username             string
	UserId               Identifier
}

type Message struct {

	MessageBody          string
	Comments             []Comment
	UploaderId           Identifier
	MessageId            Identifier
	Date                 string
	
}

type Comment struct {

	CommentId            Identifier
	MessageId            Identifier
	CommentBody          string
	Date                 string
	UploaderId           Identifier

}

type Photo struct {

	PhotoId              Identifier
	Path                 string
	UploaderId           Identifier
	Date                 string
	Comments             []Comment
}

type Group struct {

	GroupId              Identifier
	Users                []User
	Messages             []Message
	GroupName            string
	GroupPhoto           Photo
}

type Conversation struct {

	ConversationId       Identifier
	UserConnected        User
	Messages             []Message
	ChatName             string


}

