import React from "react";
import Thread from "../models/Thread";
import ThreadCard from "./ThreadCards";

interface ThreadListProps {
  threads: Thread[]
}

export default function ThreadList({ threads }: ThreadListProps) {
  return (
    <div className="thread-list">
      {threads.map(thread => <ThreadCard key={thread.id} thread={thread} />)}
    </div>
  )
}