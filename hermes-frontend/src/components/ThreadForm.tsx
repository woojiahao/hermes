import React, {createRef, useEffect, useState} from "react"
import {emptyThread, ThreadDto} from "../models/Thread"
import {ContentState, convertToRaw} from "draft-js"
import Tag from "../models/Tag"
import RichTextEditor from "./RichTextEditor"
import TagSelection from "./TagSelection"
import {draftToMarkdown, markdownToDraft} from "markdown-draft-js"
import {HermesRequest, jsonConvert} from "../utility/request"
import {asyncGetThread} from "../services/thread"
import {useNavigate} from "react-router-dom"

interface ThreadFormProps {
  threadId: string
  action: (dto: ThreadDto) => Promise<void>
  error: string
  setError: (error: string) => void
}

export default function ThreadForm({threadId, action, error, setError}: ThreadFormProps) {
  const titleRef = createRef<HTMLInputElement>()
  const [thread, setThread] = useState(emptyThread())
  const [tagsLoaded, setTagsLoaded] = useState(false)
  const [contentState, setContentState] = useState(null)
  const [selectedTags, setSelectedTags] = useState<Map<number, Tag>>(new Map())
  const navigate = useNavigate()

  useEffect(() => {
    if (!threadId) {
      setContentState(convertToRaw(ContentState.createFromText("")))
    } else {
      (async () => {
          await asyncGetThread(
            threadId,
            thread => {
              setThread(thread)
              const content = markdownToDraft(thread.content)
              setContentState(content)
            },
            e => setError(e.message),
            e => setError(e.message)
          )
        }
      )()
    }
  }, [])

  useEffect(() => {
    if (thread.tags.length > 0 && !tagsLoaded) {
      // Load the initial set of selected tags
      const initialTagContents = thread.tags.map(tag => tag.content);

      (async () => {
        await new HermesRequest()
          .GET()
          .endpoint("/threads/tags")
          .onSuccess(json => {
            const tags = jsonConvert.deserializeArray(json, Tag)
            const filteredTags = tags.filter(tag => initialTagContents.includes(tag.content))
            const map = new Map()
            filteredTags.forEach((tag, i) => {
              map.set(i, tag)
            })
            setSelectedTags(map)
            setTagsLoaded(true)
          })
          .call()
      })()
    }
  }, [thread])

  function validateInput(): [boolean, string?, string?, string?] {
    const title = titleRef.current.value.trim() as string
    const content = (draftToMarkdown(contentState) as string).trim()

    if (title.length < 5) {
      return [false, "Thread title must be at least 5 characters long"]
    }

    if (content.length < 30) {
      return [false, "Thread content must be at least 30 characters long"]
    }

    return [true, null, title, content]
  }

  async function onSubmit() {
    setError("")

    const [status, error, title, content] = validateInput()
    if (!status) {
      setError(error)
      return
    }

    const dto: ThreadDto = {
      title: title,
      content: content,
      tags:
        Array.from(selectedTags.values()).map(tag => {
          return {
            'content': tag.content,
            'hex_code': tag.hexCode
          }
        }),
      is_published: thread.isPublished,
      is_open: thread.isOpen
    }

    await action(dto)
  }

  return (
    <div className="form card thick">
      {error && <p className="error">{error}</p>}

      <div className="field">
        <p>Thread title</p>
        <span>Use an interesting and catchy title!</span>
        <input type="text" name="title" id="title" ref={titleRef} defaultValue={thread.title}/>
      </div>

      <div className="field">
        <p>Thread content</p>
        <span><em>hermes</em> supports rich text editing!</span>
        {contentState ?
          <RichTextEditor contentState={contentState} setContentState={setContentState}/> :
          <p className="error">Something happened with RichTextView</p>
        }
      </div>

      <div className="field">
        <p>Tags</p>
        <span>Select tags for this thread to be easily identified.</span>
        <TagSelection selectedTags={selectedTags} setSelectedTags={setSelectedTags}/>
      </div>

      <div className="flex justify-between">
        <button type="button" className="button" onClick={() => navigate(-1)}>Cancel</button>
        <button type="button" className="button blue" onClick={async () => await onSubmit()}>Submit
        </button>
      </div>
    </div>
  )
}
