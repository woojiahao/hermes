import React from "react";
import ReactMarkdown from "react-markdown"
import Thread from "../models/Thread";
import {useNavigate} from "react-router-dom"
import DisplayTag from "./DisplayTag"
import {formatDate} from "../utility/general"

interface ThreadCardProps {
  thread: Thread
}

export default function ThreadCard({thread}: ThreadCardProps) {
  const navigate = useNavigate()

  return (
    <div className="thin-card thread-card" onClick={() => navigate(`/threads/${thread.id}`)}>
      <h3 className="thread-title">{thread.title}</h3>
      <ReactMarkdown className="markdown preview">{thread.content}</ReactMarkdown>
      {thread.tags &&
        <div className="thread-tags" hidden={thread.tags.length === 0}>
          {thread.tags.map((tag, i) => <DisplayTag key={i} tag={tag}></DisplayTag>)}
        </div>}
      <div className="ends">
        <p className="subtitle">Posted by {thread.creator}</p>
        <p className="subtitle">Posted on {formatDate(thread.createdAt)}</p>
      </div>
    </div>
  )
}
