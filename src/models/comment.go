package models

import(
)

type(
	Comment struct {
		ID			int					`json:"id"`
		Author		User					`json:"author"`
		Content		string				`json:"content"`
		Date		string	`json:"created"`
		Likes		int					`json:"likes"`
		Liked		bool	`json:"liked"`
		PostID		int					`json:"postId"`
		IsFollowingAuthor bool `json:"isFollowing"`
	}
)