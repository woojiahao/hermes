import {JsonObject, JsonProperty} from "json2typescript"

@JsonObject("Comment")
export default class Comment {
  @JsonProperty("id", String)
  id: string

  @JsonProperty("content", String)
  content: string
}
