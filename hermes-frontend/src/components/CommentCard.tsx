import React from "react"
import Comment from "../models/Comment"
import {useAppSelector} from "../redux/hooks"
import {MdDeleteOutline} from "react-icons/md"
import {formatDate} from "../utility/general"

interface CommentCardProps {
  deleteComment: () => void
  comment: Comment
}

export default function CommentCard({deleteComment, comment}: CommentCardProps) {
  const user = useAppSelector((state) => state.user.user)

  return (
    <div className="thin-card comment-card">
      <div>
        <p className="comment-card-by subtitle">Comment by {comment.creator} on {formatDate(comment.createdAt)}</p>
        <p>{comment.content}</p>
      </div>

      {user && (comment.createdBy === user.id || user.role === 'ADMIN') &&
        <MdDeleteOutline onClick={deleteComment} size={25} color={"var(--error-primary)"}/>}
    </div>
  )
}
