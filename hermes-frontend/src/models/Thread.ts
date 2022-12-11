// TODO: Display the OP's username
import { JsonObject, JsonProperty } from "json2typescript";
import Tag from "./Tag";

@JsonObject("Thread")
export default class Thread {
  @JsonProperty("id", String)
  id: string

  @JsonProperty("is_published", Boolean)
  isPublished: boolean

  @JsonProperty("is_open", Boolean)
  isOpen: boolean

  @JsonProperty("title", String)
  title: string

  @JsonProperty("content", String)
  content: string

  @JsonProperty("tags", [Tag])
  tags: Tag[]
}