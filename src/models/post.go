package models


type(
	Post struct {
		ID			int					`json:"id"`
		Author		User					`json:"author"`
		Content		string				`json:"content"`
		Date		string	`json:"created"`
		Likes		int					`json:"likes"`
		Liked 	bool		`json:"liked"`
		Comments	[]*Comment			`json:"comments"`
		IsFollowingAuthor bool `json:"isFollowing"`
	}
)