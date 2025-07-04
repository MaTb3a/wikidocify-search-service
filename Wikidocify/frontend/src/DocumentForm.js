import React, {useState } from "react";

function DocumentForm() {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage("");
    setLoading(true);

    try {
      const res = await fetch("http://localhost:8081/documents", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, content }),
      });
      if (res.ok) {
        setMessage("Document added successfully!");
        setTitle("");
        setContent("");
      } else {
        const err = await res.text();
        setMessage("Error: " + err);
      }
    } catch (err) {
      setMessage("Error: " + err.message);
    }
    setLoading(false);
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>Add Document</h2>
      <div style={{ marginBottom: 12 }}>
        <label>
          Title:
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            style={{ width: "100%", padding: 8, marginTop: 4 }}
          />
        </label>
      </div>
      <div style={{ marginBottom: 12 }}>
        <label>
          Content:
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            required
            rows={4}
            style={{ width: "100%", padding: 8, marginTop: 4 }}
          />
        </label>
      </div>
      <button type="submit" disabled={loading}>
        {loading ? "Sending..." : "Send"}
      </button>
      {message && (
        <div
          style={{
            marginTop: 16,
            color: message.startsWith("Error") ? "red" : "green",
          }}
        >
          {message}
        </div>
      )}
    </form>
  );
}

export default DocumentForm;
