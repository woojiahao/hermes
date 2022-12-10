import { Cookies } from "react-cookie";

export function getJWT(): string | undefined {
  const cookies = new Cookies()
  return cookies.get('token')
}