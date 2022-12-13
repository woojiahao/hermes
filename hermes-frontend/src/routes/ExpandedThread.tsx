import React, {useEffect, useRef, useState} from "react"
import Thread, {emptyThread} from "../models/Thread"
import {useNavigate, useParams} from "react-router-dom"
import {errorFields, errorMessage, HermesRequest, jsonConvert} from "../utility/request"
import {IoArrowBackSharp} from "react-icons/io5"
import {useAppSelector} from "../redux/hooks"
import ReactMarkdown from "react-markdown"
import Comment from "../models/Comment"

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
        .hasAuthorization()
        .endpoint(`/threads/${id}`)
        .onSuccess(json => {
          const t = jsonConvert.deserializeObject(json, Thread)
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
      .hasAuthorization()
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

  return (
    <div className="single">
      <div className="menu">
        <div className="group">
          <IoArrowBackSharp onClick={() => navigate(-1)}/>
          <h1 className="heading">Thread</h1>
        </div>
        {thread.publisher === user.id && <a href="/edit-thread" className='button'>Edit</a>}
      </div>

      <div className="expanded-thread">
        {error && <p className="error">{error}</p>}
        <h2>{thread.title}</h2>
        <ReactMarkdown>{thread.content}</ReactMarkdown>

        <div className="comments">
          <h3>Comments</h3>
          <textarea name="new-comment" id="new-comment" placeholder="Leave a comment" cols={30} rows={10}
                    ref={commentRef}></textarea>
          <button type="button" className="button" onClick={submitComment}>Submit Comment</button>
          <div>
            {comments.map(comment => <p key={comment.id}>{comment.content}</p>)}
          </div>
        </div>
      </div>
    </div>
  )
}
