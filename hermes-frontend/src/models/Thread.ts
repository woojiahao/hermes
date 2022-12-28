import {JsonObject, JsonProperty} from "json2typescript";
import Tag from "./Tag";
import {DateConverter} from "./DateConverter"

export interface ThreadDto {
  title: string
  content: string
  tags: { 'content': string, 'hex_code': string }[]
  is_published: boolean
  is_open: boolean
}

@JsonObject("Thread")
export default class Thread {
  @JsonProperty("id", String)
  id: string

  @JsonProperty("is_published", Boolean)
  isPublished: boolean

  @JsonProperty("is_open", Boolean)
  isOpen: boolean

  @JsonProperty("is_pinned", Boolean)
  isPinned: boolean

  @JsonProperty("title", String)
  title: string

  @JsonProperty("content", String)
  content: string

  @JsonProperty("tags", [Tag])
  tags: Tag[]

  @JsonProperty("created_at", DateConverter)
  createdAt: Date

  @JsonProperty("created_by", String)
  createdBy: string

  @JsonProperty("creator", String)
  creator: string

  @JsonProperty("upvoters", [String])
  upvoters: string[]

  @JsonProperty("downvoters", [String])
  downvoters: string[]

  @JsonProperty("upvotes")
  upvotes: number

  @JsonProperty("downvotes")
  downvotes: number
}



export function emptyThread(): Thread {
  return {
    id: "",
    isPublished: false,
    isOpen: false,
    isPinned: false,
    title: "",
    content: "",
    tags: [],
    createdAt: new Date(),
    createdBy: "",
    creator: "",
    upvoters: [],
    downvoters: [],
    upvotes: 0,
    downvotes: 0,
  }
}
