import React, {createRef, useEffect, useRef, useState} from "react";
import {HermesRequest, jsonConvert} from "../utility/request"
import Tag from "../models/Tag"
import {MdOutlineCheck, MdOutlineClose} from "react-icons/md"
import {tagStyle} from "../utility/tag"
import DisplayTag from "./DisplayTag"

interface TagSelectionProps {
  selectedTags: Map<number, Tag>,
  setSelectedTags: React.Dispatch<React.SetStateAction<Map<number, Tag>>>
}

export default function TagSelection({selectedTags, setSelectedTags}: TagSelectionProps) {
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
  }, [selectedTags])

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
  }, [])

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
    setSelectedTags(prevState => {
      const cur = new Map(prevState)
      cur.set(id, selectedTag)
      return cur
    })
  }

  function removeSelection(id: number) {
    setSelectedTags(prevState => {
      const cur = new Map(prevState)
      cur.delete(id)
      return cur
    })
  }

  return (
    <div className="w-full" ref={componentRef} onFocus={() => setIsShown(true)}>
      {error && <p className="error">{error}</p>}

      <div>
        <div className="flex flex-row flex-wrap gap-2">
          {selectedTags &&
            Array
              .from(selectedTags.entries())
              .map(([i, tag]) =>
                <div key={i} style={tagStyle(tag)} className="flex justify-center items-center gap-y-2 w-fit px-1.5 py-2 rounded-br">
                  <span style={tagStyle(tag)}>{tag.content}</span><MdOutlineClose className="hover:cursor-pointer ml-2" onClick={() => removeSelection(i)}/>
                </div>
              )
          }
        </div>
        <input type="text"
               className="w-full"
               onChange={() => setSearchTerm(searchTagRef.current.value)}
               ref={searchTagRef}/>
      </div>

      <div className="border border-primary bg-background-secondary border-t-0 rounded-b-[8px]" hidden={!isShown}>
        <div className="flex flex-row items-center">
          <input
            type="text"
            name="new-tag-content"
            id="new-tag-content"
            placeholder="New Tag"
            defaultValue=""
            onKeyDown={e => {
              if (e.key === 'Enter') addNewTag()
            }}
            className="border-0 flex-1 focus:outline-0 m-0 p-2"
            ref={newTagContentRef}/>
          <input
            type="color"
            name="new-tag-color"
            id="new-tag-color"
            className="border-0 bg-background-secondary w-[30px] h-[30px] !p-0 !m-0 hover:cursor-pointer"
            ref={newTagHexCodeRef}/>
        </div>

        <div className="tags-list">
          {tags &&
            tags.filter(tag => tag.content.includes(searchTerm))
              .map((tag, i) =>
                <div key={i}
                     className="flex p-2 items-center justify-between hover:cursor-pointer hover:bg-[#e7e7e7] last:rounded-b-br"
                     onClick={() => {
                       if (!(i in selectedTags)) selectTag(i)
                     }}>
                  <DisplayTag key={i} tag={tag}></DisplayTag>
                  {Array.from(selectedTags.values()).filter(t => t.content === tag.content).length > 0 && <MdOutlineCheck size={25} />}
                </div>)}
        </div>
      </div>
    </div>
  )
}
