package enum

type CommentStatus string

const (
	NoComments  CommentStatus = "no_comments"  // 无评论状态
	HasComments CommentStatus = "has_comments" // 有评论状态
)
