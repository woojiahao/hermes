import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import './index.css';
import About from './routes/about';
import Home from './routes/home';
import Login from './routes/login';

const root = ReactDOM.createRoot(document.getElementById('root'));

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home></Home>,
  },
  {
    path: "/about",
    element: <About></About>,
  },
  {
    path: "/login",
    element: <Login></Login>,
  }
])

root.render(
  <React.StrictMode>
    <div className="container">
      <header>
        <h1>hermes</h1>
        {/* TODO: Dynamically generate navigation depending on whether they are logged in or not */}
        <nav>
          <a href="/">Home</a>
          <a href="/about">About</a>
          <a href="/login" className="login-button">Login</a>
        </nav>
      </header>
      <RouterProvider router={router} />
      <footer>
        <p>Copyright &copy; 2022 (Woo Jia Hao)</p>
        <p>hermes is a web forum designed with ❤️ using React and Go</p>
      </footer>
    </div>
  </React.StrictMode>
);