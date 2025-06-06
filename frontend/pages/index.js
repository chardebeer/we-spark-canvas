// frontend/pages/index.js
import useImages from "../hooks/useImages";
import ImageCard from "../components/ImageCard";

export default function Home() {
  const { images, loading, error } = useImages(20, 0);

  if (loading) return <p>Loading…</p>;
  if (error) return <p>Error loading images.</p>;

  return (
    <div className="max-w-5xl mx-auto p-4">
      <h1 className="text-3xl font-semibold mb-6">We Spark Canvas</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        {images.map((img) => (
          <ImageCard key={img.id} image={img} />
        ))}
      </div>
    </div>
  );
}
