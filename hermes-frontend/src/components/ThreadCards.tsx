import React from "react";
import Thread from "../models/Thread";
import {tagStyle} from "../utility/tag"

interface ThreadCardProps {
  thread: Thread
}

export default function ThreadCard({thread}: ThreadCardProps) {
  console.log(thread)
  return (
    <div className="thread-card">
      <h3 className="title">{thread.title}</h3>
      <p className="subtitle">{thread.content}</p>
      {thread.tags &&
        <div className="tags" hidden={thread.tags.length === 0}>
          {thread.tags.map((tag, i) => <span key={i} style={tagStyle(tag)}>{tag.content}</span>)}
        </div>
      }
    </div>
  )
}
