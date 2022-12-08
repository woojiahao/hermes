import React from 'react';
import ThreadList from '../components/ThreadList';
import Thread from '../models/Thread';

export default function Home() {
  const threads: Thread[] = [
    { id: "sdf", isPublished: true, isOpen: false, title: "hello world", content: "weofiqjoewijfqoi", createdAt: new Date() },
    { id: "sof", isPublished: true, isOpen: false, title: "hello world", content: "weofiqjoewijfqoi", createdAt: new Date() },
    { id: "xdf", isPublished: true, isOpen: false, title: "hello world", content: "weofiqjoewijfqoi", createdAt: new Date() },
    { id: "sqf", isPublished: true, isOpen: false, title: "hello world", content: "weofiqjoewijfqoi", createdAt: new Date() },
    { id: "slf", isPublished: true, isOpen: false, title: "hello world", content: "weofiqjoewijfqoi", createdAt: new Date() },
    { id: "soo", isPublished: true, isOpen: false, title: "hello world", content: "weofiqjoewijfqoi", createdAt: new Date() },
  ]

  return (
    <div className="content">
      <h1>Threads</h1>
      <ThreadList threads={threads}></ThreadList>
    </div>
  );
}
