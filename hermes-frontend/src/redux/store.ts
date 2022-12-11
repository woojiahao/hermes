import { configureStore } from "@reduxjs/toolkit";
import authReducer from './authSlice';
import userReducer from './userSlice';


const store = configureStore({
  middleware: (getDefaultMiddleware) => getDefaultMiddleware({
    serializableCheck: false
  }),
  reducer: {
    auth: authReducer,
    user: userReducer
  }
});

export default store

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch