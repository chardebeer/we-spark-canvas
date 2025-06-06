// frontend/pages/collections/index.js
import { useState, useEffect } from "react";
import Link from "next/link";
import apiClient from "../../lib/axios";

export default function Collections() {
  const [collections, setCollections] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    apiClient
      .get("/collections")  // You’ll need to implement GET /collections in Go if not already
      .then((res) => {
        setCollections(res.data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err);
        setLoading(false);
      });
  }, []);

  if (loading) return <p>Loading…</p>;
  if (error) return <p>Error loading collections.</p>;

  return (
    <div className="max-w-5xl mx-auto p-4">
      <h1 className="text-3xl font-semibold mb-6">Collections</h1>
      <ul className="space-y-4">
        {collections.map((col) => (
          <li key={col.id}>
            <Link href={`/collections/${col.id}`} className="text-blue-600 hover:underline">
              {col.title}
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
