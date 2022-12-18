import User from "../models/User";
import {clearJWT, hasValidJWT} from "../utility/jwt";
import {HermesRequest, jsonConvert} from "../utility/request";

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
      .onFailure(_ => {
        clearJWT()
      })
      .call()
  }

  return user
}

export async function createUser(
  username: string,
  password: string,
  role: 'ADMIN' | 'USER',
  onSuccess: (user: User) => void,
  onFailure: (error: string) => void,
  onError: (e: Error) => void
) {
  await new HermesRequest()
    .POST()
    .endpoint('register')
    .body({
      'username': username,
      'password': password,
      'role': role,
    })
    .onSuccess(json => {
      const user = jsonConvert.deserializeObject(json, User)
      onSuccess(user)
    })
    .onFailure((e) => {
      let message = ""
      if (typeof e.message === 'string') {
        message = e.message
      } else {
        message = e.message.map((e) => `${e.field} ${e.message}`).join(". ")
      }

      onFailure(message)
    })
    .onError(e => onError(e))
    .call()
}
