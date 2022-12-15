import React, {useEffect, useRef, useState} from "react"
import Thread, {emptyThread} from "../models/Thread"
import {useNavigate, useParams} from "react-router-dom"
import {errorFields, errorMessage, HermesRequest, jsonConvert} from "../utility/request"
import {IoArrowBackSharp} from "react-icons/io5"
import {useAppSelector} from "../redux/hooks"
import ReactMarkdown from "react-markdown"
import Comment from "../models/Comment"
import {formatDate} from "../utility/general"
import CommentCard from "../components/CommentCard"

export default function ExpandedThread() {
  const {id} = useParams()
  const [thread, setThread] = useState<Thread>(emptyThread())
  const [comments, setComments] = useState<Comment[]>([])
  const [error, setError] = useState<string>()
  const navigate = useNavigate()
  const user = useAppSelector((state) => state.user.user)
  const commentRef = useRef<HTMLTextAreaElement>()

  useEffect(() => {
    (async () => {
      await new HermesRequest()
        .GET()
        .endpoint(`/threads/${id}`)
        .onSuccess(json => {
          const t = jsonConvert.deserializeObject(json, Thread)
          console.log(t)
          setThread(t)
        })
        .onFailure((e: errorMessage) => setError(e.message))
        .onError(e => setError(e.message))
        .call()
      await getComments()
    })()
  }, [])

  async function getComments() {
    await new HermesRequest()
      .GET()
      .endpoint(`/threads/${id}/comments`)
      .onSuccess(json => {
        const c = jsonConvert.deserializeArray(json, Comment)
        console.log(c)
        setComments(c)
      })
      .onFailure((e: errorMessage) => setError(e.message))
      .onError(e => setError(e.message))
      .call()
  }

  async function submitComment() {
    setError('')
    await new HermesRequest()
      .POST()
      .endpoint(`/threads/${id}/comments`)
      .hasAuthorization()
      .body({
        'user_id': user.id,
        'thread_id': id,
        'content': commentRef.current.value
      })
      .onSuccess(_ => {
        getComments()
        commentRef.current.value = ""
      })
      .onFailure((e: errorMessage | errorFields) => {
        // TODO: Abstract this
        if (typeof e.message === 'string') {
          setError(e.message)
        } else {
          const formatted = e.message.map((e) => `${e.field} ${e.message}`).join(". ")
          setError(formatted)
        }
      })
      .onError(e => setError(e.message))
      .call()
  }

  async function deleteThread() {
    await new HermesRequest()
      .DELETE()
      .endpoint(`/threads/${thread.id}`)
      .hasAuthorization()
      .onSuccess(_ => {
        navigate('/')
      })
      .call()
  }

  async function deleteComment(commentId: string) {
    await new HermesRequest()
      .DELETE()
      .endpoint(`/threads/${thread.id}/comments/${commentId}`)
      .hasAuthorization()
      .onSuccess(_ => {
        getComments()
      })
      .call()
  }

  return (
    <div className="single">
      <div className="title">
        <div className="group">
          <IoArrowBackSharp onClick={() => navigate(-1)} size={25} color={`var(--primary-color)`}/>
          <h1 className="heading">Thread</h1>
        </div>
        <div className="group">
          {user && (thread.createdBy === user.id || user.role === 'ADMIN') &&
            <button type="button" onClick={deleteThread} className='static-button-red'>Delete</button>}
          {user && (thread.createdBy === user.id || user.role === 'ADMIN') &&
            <a href="/edit-thread" className='static-button-blue'>Edit</a>}
        </div>
      </div>

      <div className="expanded-thread">
        {error && <p className="error">{error}</p>}
        <h2>{thread.title}</h2>
        <div className="ends">
          <p className="subtitle">Posted by {thread && thread.creator}</p>
          <p className="subtitle">Posted by {thread && formatDate(thread.createdAt)}</p>
        </div>
        <ReactMarkdown>{thread.content}</ReactMarkdown>

        <div className="comments">
          <h3>Comments</h3>
          {user && <div>
          <textarea name="new-comment" id="new-comment" placeholder="Leave a comment" cols={30} rows={10}
                    ref={commentRef}></textarea>
            <button type="button" className="static-button-blue" onClick={submitComment}>Submit Comment</button>
          </div>}
          {comments.length > 0 ?
            <div className="comments-list">
              {comments.map(comment => <CommentCard key={comment.id}
                                                    deleteComment={async () => await deleteComment(comment.id)}
                                                    comment={comment}/>)}
            </div> :
            <p>No comments yet!</p>
          }
        </div>
      </div>
    </div>
  )
}
