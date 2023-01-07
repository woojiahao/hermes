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
import {BsFillHandThumbsDownFill, BsFillHandThumbsUpFill, BsPinAngle, BsPinAngleFill} from "react-icons/bs"

type VoteType = "UPVOTE" | "DOWNVOTE" | "NONE"

export default function ExpandedThread() {
  const {id} = useParams()
  const [thread, setThread] = useState<Thread>(emptyThread())
  const [comments, setComments] = useState<Comment[]>([])
  const [error, setError] = useState<string>()
  const [vote, setVote] = useState<VoteType>("NONE")

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

  useEffect(() => {
    if (thread && user) {
      if (thread.upvoters.includes(user.id)) setVote("UPVOTE")
      else if (thread.downvoters.includes(user.id)) setVote("DOWNVOTE")
    }
  }, [thread, user])

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
    if (commentRef.current.value.trim().length < 5) {
      setError('Comment length must be at least 5 characters long')
      return
    }

    await new HermesRequest()
      .POST()
      .endpoint(`/threads/${id}/comments`)
      .hasAuthorization()
      .body({
        'user_id': user.id,
        'thread_id': id,
        'content': commentRef.current.value.trim()
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

  async function makeVote(voteType: VoteType) {
    if (!user) {
      setError("Login before voting")
      return
    }

    if (vote === voteType) {
      // Remove the existing vote
      await new HermesRequest()
        .endpoint(`threads/${id}/votes`)
        .hasAuthorization()
        .DELETE()
        .onSuccess(_ => {
          if (vote === "UPVOTE") {
            thread.upvoters = thread.upvoters.filter(i => i !== user.id)
            thread.upvotes -= 1
          } else {
            thread.downvoters = thread.downvoters.filter(i => i !== user.id)
            thread.downvotes -= 1
          }
          setVote("NONE")
        })
        .onFailure((e: { message: string }) => setError(e.message))
        .onError(e => setError(e.message))
        .call()
    } else {
      // Change the vote to the other
      await new HermesRequest()
        .endpoint(`threads/${id}/votes`)
        .hasAuthorization()
        .PUT()
        .body({
          "is_upvote": voteType === "UPVOTE"
        })
        .onSuccess(_ => {
          if (voteType === "UPVOTE") {
            thread.upvoters.push(user.id)
            thread.upvotes += 1
            if (vote !== 'NONE') {
              thread.downvotes -= 1
            }
          } else {
            thread.downvoters.push(user.id)
            thread.downvotes += 1
            if (vote !== 'NONE') {
              thread.upvotes -= 1
            }
          }
          setVote(voteType)
        })
        .onFailure((e: { message: string }) => setError(e.message))
        .onError(e => setError(e.message))
        .call()
    }
  }

  return (
    <Layout>
      <div className="single">
        <div className="title tiny:flex tiny:flex-col tiny:gap-y-4">
          <div className="group tiny:w-full tiny:flex tiny:justify-start">
            <IoArrowBackSharp onClick={() => navigate(-1)} size={25} color={`var(--primary-color)`}/>
            <h1 className="heading">Thread</h1>
          </div>
          <div className="group tiny:w-full tiny:flex tiny:justify-between">
            {user && (thread.createdBy === user.id || user.role === 'ADMIN') &&
              <button type="button" onClick={deleteThread} className='button red'>Delete</button>}
            {user && (thread.createdBy === user.id || user.role === 'ADMIN') &&
              <a href={`/threads/${thread.id}/edit`} className='button blue'>Edit</a>}
          </div>
        </div>

        <div>
          {error && <p className="error mb-4">{error}</p>}
          <div className="mb-8 flex flex-col gap-y-2">
            <div className="flex justify-between items-center">
              <h2 className="break-words">{thread.title}</h2>
              <div className="group tiny:flex-col-reverse tiny:gap-y-4 items-center justify-center">
                <div className="flex gap-x-2 items-center">
                  <BsFillHandThumbsUpFill size={20}
                                          className="hover:cursor-pointer"
                                          onClick={() => makeVote("UPVOTE")}
                                          color={vote === "UPVOTE" ? "#eb6b1c" : "#191919"}/>
                  <span className={vote === "UPVOTE" ? "text-[#eb6b1c]" : "text-dark"}>{thread.upvotes}</span>
                  <BsFillHandThumbsDownFill size={20}
                                            className="hover:cursor-pointer"
                                            onClick={() => makeVote("DOWNVOTE")}
                                            color={vote === "DOWNVOTE" ? "#77b2e6" : "#191919"}/>
                  <span className={vote === "DOWNVOTE" ? "text-[#77b2e6]" : "text-dark"}>{thread.downvotes}</span>
                </div>
                {thread.isPinned ?
                  <BsPinAngleFill className={`${user && user.role === 'ADMIN' ? 'clickable' : ''} hover:cursor-pointer`}
                                  color="#ebc81a"
                                  size={25} onClick={async () => await pinThread()}/> :
                  <BsPinAngle className={`${user && user.role === 'ADMIN' ? 'clickable' : ''} hover:cursor-pointer`}
                              color="#ebc81a" size={25}
                              onClick={async () => await pinThread()}/>
                }
              </div>
            </div>
            <div className="flex gap-3">
              {thread.tags.map((tag, i) => <DisplayTag tag={tag} key={i}/>)}
            </div>
            <div className="flex justify-between">
              <p className="text-dark-secondary italic">Posted by {thread && thread.creator}</p>
              <p className="text-dark-secondary italic">Posted by {thread && formatDate(thread.createdAt)}</p>
            </div>
          </div>
          <ReactMarkdown className="markdown card thick">{thread.content}</ReactMarkdown>

          <div className="mt-8 text-primary">
            <h3 className="mb-4">Comments</h3>
            {user && <div>
          <textarea name="new-comment"
                    id="new-comment"
                    placeholder="Leave a comment"
                    cols={30}
                    rows={10}
                    ref={commentRef}
                    className="w-full border rounded-br border-primary p-4 font-sans text-base mb-4 text-dark"/>
              <button type="button" className="button blue mb-4" onClick={submitComment}>Submit Comment
              </button>
            </div>}
            {comments.length > 0 ?
              <div className="flex flex-col gap-y-4">
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
