import React from "react";
import ReactMarkdown from "react-markdown"
import Thread from "../models/Thread";
import {tagStyle} from "../utility/tag"
import {useNavigate} from "react-router-dom"

interface ThreadCardProps {
  thread: Thread
}

export default function ThreadCard({thread}: ThreadCardProps) {
  const navigate = useNavigate()

  return (
    <div className="thread-card" onClick={() => navigate(`/threads/${thread.id}`)}>
      <h3 className="title">{thread.title}</h3>
      <ReactMarkdown>{thread.content}</ReactMarkdown>
      {thread.tags &&
        <div className="tags" hidden={thread.tags.length === 0}>
          {thread.tags.map((tag, i) => <span key={i} style={tagStyle(tag)}>{tag.content}</span>)}
        </div>
      }
    </div>
  )
}
