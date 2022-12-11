import User from "../models/User";
import { hasValidJWT } from "../utility/jwt";
import { HermesRequest, jsonConvert } from "../utility/request";

export async function getCurrentUser(): Promise<User> {
  let user: User = null
  if (hasValidJWT()) {
    await new HermesRequest()
      .GET()
      .endpoint("/users/current")
      .hasAuthorization()
      .onSuccess((json) => {
        user = jsonConvert.deserializeObject(json, User)
      })
      .call()
  }

  return user
}