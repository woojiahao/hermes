import { JsonObject, JsonProperty } from "json2typescript"

@JsonObject("User")
export default class User {
  @JsonProperty("id", String)
  id: string

  @JsonProperty("username", String)
  username: string

  @JsonProperty("role", String)
  role: string
}