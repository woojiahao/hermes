import React from "react";
import Thread from "../models/Thread";

interface ThreadCardProps {
  thread: Thread
}

export default function ThreadCard({ thread }: ThreadCardProps) {
  return (
    <div className="thread-card">
      <h3 className="title">{thread.title}</h3>
      <p className="subtitle">{thread.content}</p>
      <p className="published-at">{thread.createdAt.toUTCString()}</p>
    </div>
  )
}