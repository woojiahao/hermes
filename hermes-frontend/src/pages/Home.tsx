import React, {useEffect, useState} from 'react';
import ThreadList from '../components/ThreadList';
import Thread from '../models/Thread';
import {errorMessage, HermesRequest, jsonConvert} from '../utility/request';
import Tag from "../models/Tag"
import DisplayTag from "../components/DisplayTag"
import Layout from "../components/Layout"
import {AiOutlineSearch} from "react-icons/ai"
import {Link} from "react-router-dom"
import ThreadSearch from "../components/ThreadSearch"

export default function Home() {
  const [threads, setThreads] = useState<Thread[]>([])
  const [message, setMessage] = useState("Loading...")
  const [tags, setTags] = useState<Tag[]>([])
  const [filteredTags, setFilteredTags] = useState<Tag[]>([])
  const [searchTerm, setSearchTerm] = useState("")

  useEffect(() => {
    (async () => {
      await new HermesRequest()
        .GET()
        .endpoint("threads")
        .onSuccess((json) => {
          const threads = jsonConvert.deserializeArray(json, Thread)
          setThreads(threads)
          setMessage("")
        })
        .onFailure((e: errorMessage) => {
          setMessage(e.message)
        })
        .onError((e: { message: string }) => {
          setMessage(e.message)
        })
        .call()

      await new HermesRequest()
        .GET()
        .endpoint("threads/tags")
        .onSuccess(json => {
          const t = jsonConvert.deserializeArray(json, Tag)
          setTags(t)
        })
        .onFailure((e: errorMessage) => setMessage(e.message))
        .onError(e => setMessage(e.message))
        .call()
    })()
  }, [])

  return (
    <Layout>
      <div className="split">
        <div className="title phone:flex-col">
          <h1 className="heading">Threads</h1>
          <div className="flex items-center gap-4 phone:justify-between phone:w-full">
            <ThreadSearch setSearchTerm={setSearchTerm}/>
            <Link to="/threads/new" className="button blue effect">New Thread</Link>
          </div>
        </div>

        <div className="content">
          <main>
            {message && <p>{message}</p>}
            {threads.length > 0 ?
              <ThreadList threads={threads.filter(thread => {
                if (!thread.title.includes(searchTerm))
                  return false

                for (const tag of filteredTags) {
                  if (!thread.tags.flatMap(t => t.content).includes(tag.content)) return false
                }
                return true
              })}/> :
              <p>No threads created yet.</p>
            }
          </main>
          <aside className="flex flex-col gap-y-4 phone:gap-y-1">
            <p><strong>Tags</strong> (click to filter)</p>
            <div hidden={filteredTags.length === 0} className="phone:flex phone:flex-row">
              {filteredTags.length > 0 && <p className="mr-4"><em>Filtered tags</em></p>}
              <div className="phone:flex phone:flex-row phone:gap-4">
                {filteredTags
                  .map((tag, i) => <DisplayTag key={i}
                                               tag={tag}
                                               onClick={() => setFilteredTags(prevState => prevState.filter(t => t !== tag))}
                                               style={{fontWeight: "bold"}}/>)}
              </div>
            </div>
            <div className="phone:flex phone:flex-row phone:gap-4">
              {filteredTags.length > 0 && <p className="mr-4"><em>Other tags</em></p>}
              {tags
                .filter(tag => !(filteredTags.includes(tag)))
                .map((tag, i) => <DisplayTag key={i}
                                             tag={tag}
                                             onClick={() => setFilteredTags(prevState => [...prevState, tag])}/>)}
            </div>
          </aside>
        </div>
      </div>
    </Layout>
  )
}
