// frontend/hooks/useImages.js
import { useState, useEffect } from "react";
import apiClient from "../lib/axios";

export default function useImages(limit = 20, offset = 0) {
  const [images, setImages] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    apiClient
      .get(`/images?limit=${limit}&offset=${offset}`)
      .then((res) => {
        setImages(res.data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err);
        setLoading(false);
      });
  }, [limit, offset]);

  return { images, loading, error };
}
