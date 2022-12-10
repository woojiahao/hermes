import { createSlice } from "@reduxjs/toolkit";
import { Cookies } from "react-cookie";
import { RootState } from "./store";

interface AuthState {
  value: boolean
}

const initialState: AuthState = {
  value: new Cookies().get('token') !== undefined
}

export const authSlice = createSlice({
  name: 'isLoggedIn',
  initialState,
  reducers: {
    toggle: (state) => {
      state.value = !state.value
    }
  }
})

export const { toggle } = authSlice.actions

export const selectAuth = (state: RootState) => state.auth.value

export default authSlice.reducer