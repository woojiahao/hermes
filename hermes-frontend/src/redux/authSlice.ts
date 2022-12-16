import {createAsyncThunk, createSlice} from "@reduxjs/toolkit";
import {Cookies} from "react-cookie";
import {getCurrentUser} from "../services/user"
import User from "../models/User"

export type UserState = "loading" | "unauthorized" | "authorized"

interface AuthState {
  isLoggedIn: boolean
  user?: User
  userState: UserState
}

const initialState: AuthState = {
  isLoggedIn: new Cookies().get('token') !== undefined,
  user: null,
  userState: "loading"
}

export const login = createAsyncThunk(
  'auth/login',
  async () => {
    let user = null
    const token = new Cookies().get('token')
    if (token) {
      user = await getCurrentUser()
    }

    return user
  }
)

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    logout: state => {
      state.isLoggedIn = false
      state.user = null
      state.userState = "loading"
    },
  },
  extraReducers: builder => {
    builder.addCase(login.pending, state => {
      state.isLoggedIn = false
      state.user = null
      state.userState = "loading"
    })
    builder.addCase(login.fulfilled, (state, action) => {
      if (action.payload) {
        state.isLoggedIn = true
        state.user = action.payload
        state.userState = "authorized"
      } else {
        state.isLoggedIn = false
        state.user = null
        state.userState = "unauthorized"
      }
    })
  }
})

export const {logout} = authSlice.actions

export default authSlice.reducer
