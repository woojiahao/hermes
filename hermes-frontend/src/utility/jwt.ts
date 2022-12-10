import { Cookies } from "react-cookie";
import store from "../redux/store";
import { toggle } from "../redux/authSlice";
import { HermesRequest } from "./request";

const cookies = new Cookies()

const jwtCookiesKey = 'token'

export type jwtToken = string

export function setJWT(jwt: jwtToken) {
  cookies.set(jwtCookiesKey, jwt)
}

export function getJWT(): jwtToken | undefined {
  return cookies.get(jwtCookiesKey)
}

export function hasValidJWT(): boolean {
  const jwtToken = getJWT()
  return jwtToken !== undefined
}

export function clearJWT() {
  cookies.remove(jwtCookiesKey)
}

export async function refreshJWT(): Promise<boolean> {
  let refreshed = false
  await new HermesRequest()
    .GET()
    .endpoint("refresh")
    .hasAuthorization()
    .onSuccess((s: { token: string }) => {
      setJWT(s.token)
      refreshed = true
    })
    .onFailure((_) => {
      clearJWT()
      store.dispatch(toggle())
      refreshed = false
    })
    .onError((_) => {
      clearJWT()
      store.dispatch(toggle())
      refreshed = false
    })
    .call()

  return refreshed
}