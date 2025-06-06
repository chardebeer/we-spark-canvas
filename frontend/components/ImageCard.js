// frontend/components/ImageCard.js
export default function ImageCard({ image }) {
  // If the URL already uses ipfs.io, you could replace it with your local gateway:
  const displayUrl = image.url.startsWith("ipfs://")
    ? `${process.env.NEXT_PUBLIC_IPFS_URL}/ipfs/${image.url.slice(7)}`
    : image.url.startsWith("https://ipfs.io/ipfs/")
    ? image.url.replace("https://ipfs.io/ipfs/", `${process.env.NEXT_PUBLIC_IPFS_URL}/ipfs/`)
    : image.url;

  return (
    <div className="bg-white rounded-lg overflow-hidden shadow hover:shadow-lg transition">
      <img
        src={displayUrl}
        alt={image.caption || "Image"}
        className="w-full h-48 object-cover"
      />
      <div className="p-3">
        {image.caption && (
          <p className="text-gray-700 text-sm mb-1">{image.caption}</p>
        )}
        <p className="text-gray-500 text-xs">{image.hearts} hearts</p>
      </div>
    </div>
  );
}
