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
  const user = useAppSelector((state) => state.auth.user)

  return (
    <div className="card flex justify-between break-words items-center">
      <div>
        <p className="text-dark-secondary italic mb-2">Comment by {comment.creator} on {formatDate(comment.createdAt)}</p>
        <p>{comment.content}</p>
      </div>

      {user && (comment.createdBy === user.id || user.role === 'ADMIN') &&
        <MdDeleteOutline onClick={deleteComment} size={25} className="text-error hover:cursor-pointer"/>}
    </div>
  )
}
