package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	rt.router.PUT("/groups/:groupId/photo", rt.wrap(rt.SETGROUPPHOTO))
	rt.router.PUT("/users/:userId/photo", rt.wrap(rt.SETMYPHOTO))

	rt.router.PUT("/groups/:groupId", rt.wrap(rt.SETGROUPNAME))
	rt.router.DELETE("/users/:userId/group/:groupId", rt.wrap(rt.LEAVEGROUP))
	rt.router.PUT("/users/:userId/group/:groupId", rt.wrap(rt.ADDTOGROUP))

	rt.router.POST("/users/:userId/conversations/:conversationId/messages", rt.wrap(rt.SENDMESSAGE))
	rt.router.DELETE("/users/:userId/conversations/:conversationId/messages/:messageId", rt.wrap(rt.DELETEMESSAGE))
	rt.router.PUT("/users/:userId/conversations/:conversationId/messages/:messageId", rt.wrap(rt.FORWARDMESSAGE))

	rt.router.POST("/users/:userId/messages/:messageId/comments", rt.wrap(rt.COMMENTMESSAGE))
	rt.router.DELETE("/users/:userId/messages/:messageId/comments/:commentId", rt.wrap(rt.UNCOMMENTMESSAGE))

	rt.router.GET("/users/:userId/conversations", rt.wrap(rt.GETMYCONVERSATIONS))
	rt.router.GET("/users/:userId/conversations/:conversationId", rt.wrap(rt.GETCONVERSATION))
	rt.router.POST("/users/:userId/conversations/:conversationId", rt.wrap(rt.CREATECONVERSATION))

	rt.router.POST("/login", rt.wrap(rt.DOLOGIN))
	rt.router.PUT("/users/:userId", rt.wrap(rt.SETMYUSERNAME))
	rt.router.GET("/users/:userId", rt.wrap(rt.GETUSER))

	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
