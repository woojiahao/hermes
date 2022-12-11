import React, {createRef, useEffect, useRef, useState} from "react";
import {HermesRequest, jsonConvert} from "../utility/request"
import Tag from "../models/Tag"
import {MdOutlineCheck, MdOutlineClose} from "react-icons/md"
import {tagStyle} from "../utility/tag"

export default function TagSelection(
  props: {
    selectedTags: Map<number, Tag>,
    setSelectedTags: React.Dispatch<React.SetStateAction<Map<number, Tag>>>
  }
) {
  const componentRef = useRef(null)
  const searchTagRef = createRef<HTMLInputElement>()
  const newTagContentRef = createRef<HTMLInputElement>()
  const newTagHexCodeRef = createRef<HTMLInputElement>()

  const [tags, setTags] = useState<Tag[]>([])
  const [searchTerm, setSearchTerm] = useState<string>("")
  const [isShown, setIsShown] = useState(false)
  const [error, setError] = useState("")

  useEffect(() => {
    (async () => {
      await new HermesRequest()
        .GET()
        .endpoint("threads/tags")
        .onSuccess((json) => {
          if (json) {
            const tags = jsonConvert.deserializeArray(json, Tag)
            setTags(tags)
          }
        })
        .call()
    })()
  }, [tags])

  useEffect(() => {
    const handleClickOutside = (event: Event) => {
      if (componentRef.current && !componentRef.current.contains(event.target)) {
        setIsShown(false)
      }
    }

    document.addEventListener('click', handleClickOutside, true)
    return () => {
      document.removeEventListener('click', handleClickOutside, true)
    }
  })

  function addNewTag() {
    setError("")

    const content = newTagContentRef.current.value

    // Check if content is repeated
    if (tags.filter(tag => tag.content === content).length > 0) {
      setError("Tag content must be unique")
      return
    }

    const hexCode = newTagHexCodeRef.current.value
    const newTag = new Tag()
    newTag.content = content
    newTag.hexCode = hexCode
    setTags(prevState => [...prevState, newTag])
    newTagContentRef.current.value = ""
  }

  function selectTag(id: number) {
    const selectedTag = tags[id]
    props.setSelectedTags(prevState => {
      const cur = new Map(prevState)
      cur.set(id, selectedTag)
      return cur
    })
  }

  function removeSelection(id: number) {
    props.setSelectedTags(prevState => {
      const cur = new Map(prevState)
      cur.delete(id)
      return cur
    })
  }

  return (
    <div className="tag-select" ref={componentRef} onFocus={() => setIsShown(true)}>
      {error && <p className="error">{error}</p>}

      <div className="tags-input">
        <div className="tags-selected">
          {props.selectedTags &&
            Array
              .from(props.selectedTags.entries())
              .map(([i, tag]) =>
                <div key={i} style={tagStyle(tag)}>
                  <span style={tagStyle(tag)}>{tag.content}</span><MdOutlineClose onClick={() => removeSelection(i)}/>
                </div>
              )
          }
        </div>
        <input type="text"
               className="search-tags"
               onChange={() => setSearchTerm(searchTagRef.current.value)}
               ref={searchTagRef}/>
      </div>

      <div className="tags-dropdown" hidden={!isShown}>
        <div className="add-new-tag">
          <input
            type="text"
            name="new-tag-content"
            id="new-tag-content"
            placeholder="New Tag"
            defaultValue=""
            onKeyDown={e => {
              if (e.key === 'Enter') addNewTag()
            }}
            ref={newTagContentRef}/>
          <input
            type="color"
            name="new-tag-color"
            id="new-tag-color"
            ref={newTagHexCodeRef}/>
        </div>

        <div className="tags-list">
          {tags &&
            tags.filter(tag => tag.content.includes(searchTerm))
              .map((tag, i) =>
                <div key={i}
                     onClick={() => {
                       if (!(i in props.selectedTags)) selectTag(i)
                     }}>
                  <p style={tagStyle(tag)}>{tag.content}</p>
                  {Array.from(props.selectedTags.values()).filter(t => t === tag).length > 0 && <MdOutlineCheck/>}
                </div>)}
        </div>
      </div>
    </div>
  )
}
