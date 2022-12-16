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
import {asyncGetThread} from "../services/thread"
import DisplayTag from "../components/DisplayTag"
import Layout from "../components/Layout"
import {BsPinAngle, BsPinAngleFill} from "react-icons/bs"

export default function ExpandedThread() {
  const {id} = useParams()
  const [thread, setThread] = useState<Thread>(emptyThread())
  const [comments, setComments] = useState<Comment[]>([])
  const [error, setError] = useState<string>()
  const navigate = useNavigate()
  const user = useAppSelector((state) => state.auth.user)
  const commentRef = useRef<HTMLTextAreaElement>()

  useEffect(() => {
    (async () => {
      await asyncGetThread(
        id,
        thread => setThread(thread),
        e => setError(e.message),
        e => setError(e.message)
      )
      await getComments()
    })()
  }, [])

  async function getComments() {
    await new HermesRequest()
      .GET()
      .endpoint(`/threads/${id}/comments`)
      .onSuccess(json => {
        const c = jsonConvert.deserializeArray(json, Comment)
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

  async function pinThread() {
    if (user.role !== 'ADMIN') return

    const status = !thread.isPinned
    await new HermesRequest()
      .PUT()
      .endpoint(`/threads/${thread.id}/pin`)
      .hasAuthorization()
      .body({
        "is_pinned": status,
      })
      .onSuccess(json => {
        const t = jsonConvert.deserializeObject(json, Thread)
        setThread(prevState => {
          // TODO: Implement proper returning of tags
          t.tags = prevState.tags
          t.creator = prevState.creator
          return t
        })
      })
      .onFailure((e: errorMessage) => setError(e.message))
      .onError(e => setError(e.message))
      .call()
  }

  return (
    <Layout>
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
              <a href={`/threads/${thread.id}/edit`} className='static-button-blue'>Edit</a>}
          </div>
        </div>

        <div className="expanded-thread">
          {error && <p className="error">{error}</p>}
          <div className="expanded-thread-heading">
            <div className="ends">
              <h2>{thread.title}</h2>
              {thread.isPinned ?
                <BsPinAngleFill className={`${user && user.role === 'ADMIN' ? 'clickable' : ''}`} color="#ebc81a" size={25} onClick={async () => await pinThread()}/> :
                <BsPinAngle className={`${user && user.role === 'ADMIN' ? 'clickable' : ''}`} color="#ebc81a" size={25} onClick={async () => await pinThread()}/>
              }
            </div>
            <div className="expanded-thread-tags">
              {thread.tags.map((tag, i) => <DisplayTag tag={tag} key={i}/>)}
            </div>
            <div className="ends">
              <p className="subtitle">Posted by {thread && thread.creator}</p>
              <p className="subtitle">Posted by {thread && formatDate(thread.createdAt)}</p>
            </div>
          </div>
          <ReactMarkdown className="markdown thick-card">{thread.content}</ReactMarkdown>

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
    </Layout>
  )
}
