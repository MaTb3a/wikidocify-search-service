import React, { useState } from "react";

function FileUpload() {
  const [title, setTitle] = useState("");
  const [author, setAuthor] = useState("");
  const [file, setFile] = useState(null);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState(""); // "success" or "error"

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage("");
    setMessageType("");

    if (!file) {
      setMessage("Please select a file.");
      setMessageType("error");
      return;
    }

    setLoading(true);

    const reader = new FileReader();
    reader.onload = async () => {
      const base64Content = reader.result.split(",")[1];
      const payload = {
        title: title.trim(),
        author: author.trim(),
        content: base64Content,
      };

      try {
        const res = await fetch("http://localhost:8081/documents", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });
        if (res.ok) {
          setMessage("File uploaded successfully!");
          setMessageType("success");
          setTitle("");
          setAuthor("");
          setFile(null);
        } else {
          let err = "";
          try {
            err = await res.text();
          } catch {
            err = "Unknown error";
          }
          setMessage("Upload failed: " + err);
          setMessageType("error");
        }
      } catch (err) {
        setMessage("Error: " + err.message);
        setMessageType("error");
      }
      setLoading(false);
    };
    reader.readAsDataURL(file);
  };

  return (
    <div
      style={{
        maxWidth: 480,
        margin: "2rem auto",
        padding: "2.5rem 2rem",
        borderRadius: "16px",
        boxShadow: "0 4px 24px rgba(0,0,0,0.10)",
        background: "#fff",
      }}
    >
      <div style={{ textAlign: "center", marginBottom: 28 }}>
        <img
          src="https://raw.githubusercontent.com/hossamhakim/wikidocify-assets/main/logo.png"
          alt="Wikidocify"
          style={{ height: 60, marginBottom: 8 }}
        />
        <h1 style={{ margin: 0, fontSize: 36, fontWeight: 700 }}>Wikidocify</h1>
        <p style={{ color: "#555", margin: "10px 0 0 0", fontSize: 18 }}>
          Upload and share your documents securely and instantly.
        </p>
      </div>
      <h2
        style={{
          marginTop: 0,
          fontSize: 26,
          fontWeight: 600,
          marginBottom: 24,
        }}
      >
        Upload a Document
      </h2>
      <form onSubmit={handleSubmit} autoComplete="off">
        <div style={{ marginBottom: 18 }}>
          <label style={{ display: "block", fontWeight: 500, marginBottom: 6 }}>
            Title:
          </label>
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            style={{
              width: "100%",
              padding: 10,
              borderRadius: 6,
              border: "1px solid #bbb",
              fontSize: 16,
            }}
            placeholder="Document title"
            autoFocus
          />
        </div>
        <div style={{ marginBottom: 18 }}>
          <label style={{ display: "block", fontWeight: 500, marginBottom: 6 }}>
            Author:
          </label>
          <input
            value={author}
            onChange={(e) => setAuthor(e.target.value)}
            required
            style={{
              width: "100%",
              padding: 10,
              borderRadius: 6,
              border: "1px solid #bbb",
              fontSize: 16,
            }}
            placeholder="Your name"
          />
        </div>
        <div style={{ marginBottom: 22 }}>
          <label style={{ display: "block", fontWeight: 500, marginBottom: 6 }}>
            File:
          </label>
          <input
            type="file"
            onChange={(e) => setFile(e.target.files[0])}
            required
            style={{
              marginTop: 4,
              fontSize: 16,
            }}
            accept=".pdf,.doc,.docx,.txt,.md,.png,.jpg,.jpeg"
          />
          {file && (
            <div style={{ fontSize: 14, color: "#555", marginTop: 4 }}>
              Selected: {file.name}
            </div>
          )}
        </div>
        <button
          type="submit"
          disabled={loading}
          style={{
            width: "100%",
            padding: "12px 0",
            borderRadius: 6,
            border: "none",
            background: "#1976d2",
            color: "#fff",
            fontWeight: 700,
            fontSize: 18,
            cursor: loading ? "not-allowed" : "pointer",
            boxShadow: "0 2px 8px rgba(25, 118, 210, 0.08)",
            transition: "background 0.2s",
          }}
        >
          {loading ? "Uploading..." : "Upload"}
        </button>
      </form>
      {message && (
        <div
          style={{
            marginTop: 22,
            color: messageType === "success" ? "#388e3c" : "#d32f2f",
            fontWeight: 600,
            textAlign: "center",
            fontSize: 17,
          }}
        >
          {message}
        </div>
      )}
    </div>
  );
}

export default FileUpload;
