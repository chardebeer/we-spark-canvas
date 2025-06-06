// frontend/pages/upload.js
import { useState } from "react";
import axios from "axios";

export default function Upload() {
  const [file, setFile] = useState(null);
  const [caption, setCaption] = useState("");
  const [tags, setTags] = useState("");
  const [uploadedBy, setUploadedBy] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!file || !uploadedBy) {
      alert("File and uploaded_by are required");
      return;
    }
    const formData = new FormData();
    formData.append("file", file);
    formData.append("caption", caption);
    formData.append("tags", tags);
    formData.append("uploaded_by", uploadedBy);

    try {
      const res = await axios.post(
        `${process.env.NEXT_PUBLIC_API_URL}/images`,
        formData,
        {
          headers: { "Content-Type": "multipart/form-data" },
        }
      );
      alert("Uploaded! ID: " + res.data.id);
    } catch (err) {
      console.error(err);
      alert("Upload failed.");
    }
  };

  return (
    <div className="max-w-xl mx-auto p-4">
      <h1 className="text-2xl font-semibold mb-4">Upload New Image</h1>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium">File</label>
          <input
            type="file"
            accept="image/*"
            onChange={(e) => setFile(e.target.files[0])}
            className="mt-1 block w-full"
          />
        </div>
        <div>
          <label className="block text-sm font-medium">Caption</label>
          <input
            type="text"
            value={caption}
            onChange={(e) => setCaption(e.target.value)}
            className="mt-1 block w-full border rounded p-2"
          />
        </div>
        <div>
          <label className="block text-sm font-medium">
            Tags (comma-separated)
          </label>
          <input
            type="text"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
            className="mt-1 block w-full border rounded p-2"
          />
        </div>
        <div>
          <label className="block text-sm font-medium">Uploaded By (user ID)</label>
          <input
            type="number"
            value={uploadedBy}
            onChange={(e) => setUploadedBy(e.target.value)}
            className="mt-1 block w-full border rounded p-2"
          />
        </div>
        <button
          type="submit"
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
        >
          Upload
        </button>
      </form>
    </div>
  );
}
