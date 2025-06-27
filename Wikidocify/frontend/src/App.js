import React from "react";
import "./App.css";
import FileUpload from "./FileUpload";
import DocumentList from "./DocumentList";

function App() {
  return (
    <div
      style={{ background: "#f5f6fa", minHeight: "100vh", padding: "2rem 0" }}
    >
      <FileUpload />
      <DocumentList />
    </div>
  );
}

export default App;
