import React, { useState } from "react";
import axios from "axios";

function App() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSearch = async (e) => {
    e.preventDefault();
    if (!query.trim()) {
      setError("Please enter a search query");
      return;
    }

    setLoading(true);
    setError("");
    try {
      const response = await axios.get(
        "http://localhost:8080/api/v1/search",
        { params: { query } }
      );
      setResults(response.data.documents || []);
    } catch (error) {
      console.error("Search error:", error);
      setError("Error searching: " + (error.response?.data?.error || error.message));
      setResults([]);
    }
    setLoading(false);
  };

  // Function to render content with proper newlines
  const renderContent = (content) => {
    if (!content) return "No content available";

    // Handle both \n and /n newline characters
    let processedContent = content;

    // First, convert /n to \n for consistency
    processedContent = processedContent.replace(/\/n/g, '\n');

    // Then split by \n and join with line breaks
    return processedContent.split('\n').map((line, index) => (
      <React.Fragment key={index}>
        {line}
        {index < processedContent.split('\n').length - 1 && <br />}
      </React.Fragment>
    ));
  };

  return (
    <div style={{ maxWidth: 800, margin: "auto", padding: 20 }}>
      <h2>Elasticsearch Document Search</h2>
      <form onSubmit={handleSearch}>
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Enter search string"
          style={{ width: "70%", padding: 8, fontSize: 16 }}
        />
        <button
          type="submit"
          style={{ padding: 8, marginLeft: 8, fontSize: 16 }}
          disabled={loading}
        >
          {loading ? "Searching..." : "Search"}
        </button>
      </form>

      {error && (
        <div style={{ color: "red", marginTop: 10, padding: 10, backgroundColor: "#ffe6e6", borderRadius: 4 }}>
          {error}
        </div>
      )}

      {loading && <p>Loading...</p>}

      {!loading && results.length > 0 && (
        <div style={{ marginTop: 20 }}>
          <h3>Search Results ({results.length})</h3>
          {results.map((doc, index) => (
            <div
              key={doc.id || index}
              style={{
                border: "1px solid #ddd",
                margin: "10px 0",
                padding: 15,
                borderRadius: 5,
                backgroundColor: "#f9f9f9"
              }}
            >
              <h4 style={{ margin: "0 0 10px 0", color: "#333" }}>
                {doc.title || "Untitled Document"}
              </h4>
              <p style={{ margin: "5px 0", color: "#666" }}>
                <strong>Author:</strong> {doc.author || "Unknown"}
              </p>
              <p style={{ margin: "5px 0", color: "#666" }}>
                <strong>ID:</strong> {doc.id}
              </p>
              <p style={{ margin: "5px 0", color: "#666" }}>
                <strong>Created:</strong> {doc.created_at ? new Date(doc.created_at).toLocaleDateString() : "Unknown"}
              </p>
              <div style={{ marginTop: 10 }}>
                <strong>Content:</strong>
                <div style={{
                  margin: "5px 0",
                  padding: 10,
                  backgroundColor: "white",
                  borderRadius: 3,
                  border: "1px solid #eee",
                  maxHeight: "200px",
                  overflow: "auto",
                  whiteSpace: "pre-wrap",
                  fontFamily: "monospace",
                  fontSize: "14px",
                  lineHeight: "1.4"
                }}>
                  {renderContent(doc.content)}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {!loading && results.length === 0 && query && !error && (
        <div style={{ marginTop: 20, textAlign: "center", color: "#666" }}>
          <p>No documents found for "{query}"</p>
          <p>Try a different search term or check if documents have been indexed.</p>
        </div>
      )}
    </div>
  );
}

export default App;
