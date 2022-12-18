import React from "react";
import Thread from "../models/Thread";
import ThreadCard from "./ThreadCard";

interface ThreadListProps {
  threads: Thread[]
}

export default function ThreadList({threads}: ThreadListProps) {
  return (
    <div className="flex flex-col gap-y-4">
      {threads
        .map(thread => <ThreadCard key={thread.id} thread={thread}/>)}
    </div>
  )
}
