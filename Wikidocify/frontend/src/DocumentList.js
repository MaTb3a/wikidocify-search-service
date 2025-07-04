import React, { useEffect, useState } from "react";

function DocumentList() {
  const [documents, setDocuments] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("http://localhost:8081/documents")
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch documents");
        return res.json();
      })
      .then((data) => {
        setDocuments(data.documents || []);
        setError("");
      })
      .catch((err) => {
        setError(err.message);
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div style={{ marginTop: 20 }}>Loading documents...</div>;
  }

  if (error) {
    return <div style={{ color: "red", marginTop: 20 }}>Error: {error}</div>;
  }

  return (
    <div style={{ marginTop: 40 }}>
      <h3>Uploaded Documents titles</h3>
      <ul>
        {documents.map((doc) => (
          <li key={doc.id}>{doc.title}</li>
        ))}
      </ul>
    </div>
  );
}

export default DocumentList;
