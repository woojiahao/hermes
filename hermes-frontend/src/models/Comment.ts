import {JsonObject, JsonProperty} from "json2typescript"
import {DateConverter} from "./DateConverter"

@JsonObject("Comment")
export default class Comment {
  @JsonProperty("id", String)
  id: string

  @JsonProperty("content", String)
  content: string

  @JsonProperty("created_at", DateConverter)
  createdAt: Date

  @JsonProperty("created_by", String)
  createdBy: string

  @JsonProperty("creator", String)
  creator: string
}
