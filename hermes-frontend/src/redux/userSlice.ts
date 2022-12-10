import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import User from "../models/User";
import { RootState } from "./store";

interface UserState {
  user: User
}

const initialState: UserState = {
  user: { id: "", username: "", role: "" }
}

export const userSlice = createSlice({
  name: 'currentUser',
  initialState,
  reducers: {
    load: (state: UserState, action: PayloadAction<User>) => {
      state.user = action.payload
    }
  }
})

export const { load } = userSlice.actions

export const selectUser = (state: RootState) => state.user.user

export default userSlice.reducer