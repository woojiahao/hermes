import { JsonObject, JsonProperty } from "json2typescript";

@JsonObject("Tag")
export default class Tag {
  @JsonProperty("content", String)
  content: string

  @JsonProperty("hex_code", String)
  hexCode: string
}
