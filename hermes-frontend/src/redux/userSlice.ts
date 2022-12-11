import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import User from "../models/User";
import { getCurrentUser } from "../services/user";
import { RootState } from "./store";

interface UserState {
  user: User
}

const initialState: UserState = {
  user: { id: "", username: "", role: "" }
}

export const loadCurrentUser = createAsyncThunk(
  'user/getCurrentUser',
  async () => {
    return await getCurrentUser()
  }
)

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    load: (state, payload: PayloadAction<User>) => {
      state.user = payload.payload
    }
  },
  extraReducers: (builder) => {
    builder.addCase(loadCurrentUser.fulfilled, (state, action) => {
      state.user = action.payload
    })
  }
})

export const { load } = userSlice.actions

export const selectUser = (state: RootState) => state.user.user

export default userSlice.reducer
