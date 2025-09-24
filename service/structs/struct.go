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
	UploaderUserid       Identifier
	UploaderUsername     string
}

type Comment struct {
	UploaderUserid       Identifier
	UploaderUsername     string
	MessageId            Identifier
	CommentBody          string
}
type Photo struct {
	PhotoId              Identifier
	Path                 string
}