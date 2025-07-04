import React from "react";
import "./App.css";
import DocumentForm from "./DocumentForm";
import DocumentList from "./DocumentList";
function App() {
  return (
    <div style={{ maxWidth: 600, margin: "2rem auto", padding: 24 }}>
      <h1>WikiDocify Document Service</h1>
      <DocumentForm />
      <hr style={{ margin: "2rem 0" }} />
      <DocumentList />
    </div>
  );
}

export default App;
