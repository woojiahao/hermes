import React, {useEffect} from "react";
import {CookiesProvider} from "react-cookie";
import {BrowserRouter, Route, Routes} from "react-router-dom";
import About from "./pages/About";
import CreateThread from "./pages/CreateThread";
import Home from "./pages/Home";
import Login from "./pages/Login";
import UserThreads from "./pages/UserThreads";
import ExpandedThread from "./pages/ExpandedThread"
import EditThread from "./pages/EditThread"
import CreateAdmin from "./pages/CreateAdmin"
import Error from "./pages/404"
import {login} from "./redux/authSlice"
import {useAppDispatch} from "./redux/hooks"

export default function App() {
  const dispatch = useAppDispatch()

  useEffect(() => {
    dispatch(login())
  }, [])

  return (
    <CookiesProvider>
      <BrowserRouter>
        <Routes>
          <Route index element={<Home/>}/>
          <Route path="/about" element={<About/>}/>
          <Route path="/login" element={<Login/>}/>
          <Route path="/threads/you" element={<UserThreads/>}/>
          <Route path="/threads/new" element={<CreateThread/>}/>
          <Route path="/threads/:id" element={<ExpandedThread/>}/>
          <Route path="/threads/:id/edit" element={<EditThread/>}/>
          <Route path="/admin/new" element={<CreateAdmin/>}/>
          <Route path="*" element={<Error message=""/>}/>
        </Routes>
      </BrowserRouter>
    </CookiesProvider>
  )
}
