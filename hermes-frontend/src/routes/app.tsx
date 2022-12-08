import React from 'react';

export default function App() {
  return (
    <div className="container">
      <header>
        <h3>hermes</h3>
        {/* TODO: Dynamically generate navigation depending on whether they are logged in or not */}
        <nav>
          <a href="/">Home</a>
          <a href="/about">About</a>
          <a href="/login">Login</a>
        </nav>
      </header>
      <div className="content">
      </div>
      <footer>
        <p>Copyright &copy; 2022 (Woo Jia Hao)</p>
        <p>hermes is a web forum designed with ❤️ using React and Go</p>
      </footer>
    </div>
  );
}
