// frontend/components/Navbar.js
import Link from "next/link";

export default function Navbar() {
  return (
    <nav className="bg-white shadow p-4 mb-6">
      <div className="container mx-auto flex justify-between items-center">
        <Link href="/" className="text-2xl font-bold">
         We Spark Canvas
        </Link>
        <div className="space-x-4">
          <Link href="/upload" className="text-blue-600 hover:underline">
            Upload
          </Link>
          <Link href="/collections" className="text-blue-600 hover:underline">
            Collections
          </Link>
        </div>
      </div>
    </nav>
  );
}
