// frontend/pages/users/[id].js
import { useState, useEffect } from "react";
import { useRouter } from "next/router";
import apiClient from "../../lib/axios";
import ImageCard from "../../components/ImageCard";

export default function UserProfile() {
  const router = useRouter();
  const { id } = router.query;

  const [images, setImages] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!id) return;
    apiClient
      .get(`/users/${id}/images`)
      .then((res) => {
        setImages(res.data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err);
        setLoading(false);
      });
  }, [id]);

  if (loading) return <p>Loading…</p>;
  if (error) return <p>Error loading user images.</p>;

  return (
    <div className="max-w-5xl mx-auto p-4">
      <h1 className="text-3xl font-semibold mb-6">User’s Images</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        {images.map((img) => (
          <ImageCard key={img.id} image={img} />
        ))}
      </div>
    </div>
  );
}
