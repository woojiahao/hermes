import React from "react";
import ReactMarkdown from "react-markdown"
import Thread from "../models/Thread";
import {useNavigate} from "react-router-dom"
import DisplayTag from "./DisplayTag"
import {formatDate} from "../utility/general"
import {BsPinAngleFill} from "react-icons/bs"

interface ThreadCardProps {
  thread: Thread
}

export default function ThreadCard({thread}: ThreadCardProps) {
  const navigate = useNavigate()

  return (
    <div className="card animated" onClick={() => navigate(`/threads/${thread.id}`)}>
      <div className="flex justify-between">
        <h3 className="mb-4 text-dark-highlight">{thread.title}</h3>
        {thread.isPinned && <BsPinAngleFill color="#ebc81a" size={25}/>}
      </div>
      <ReactMarkdown className="markdown preview">{thread.content}</ReactMarkdown>
      {thread.tags &&
        <div className="flex gap-3 flex-wrap" hidden={thread.tags.length === 0}>
          {thread.tags.map((tag, i) => <DisplayTag key={i} tag={tag}></DisplayTag>)}
        </div>}
      <div className="flex justify-between mt-4">
        <p className="text-dark-highlight italic">Posted by {thread.creator}</p>
        <p className="text-dark-highlight italic">Posted on {formatDate(thread.createdAt)}</p>
      </div>
    </div>
  )
}
