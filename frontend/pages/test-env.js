export default function TestEnv() {
  return (
    <div>
      <p>API URL: {process.env.NEXT_PUBLIC_API_URL}</p>
      <p>IPFS URL: {process.env.NEXT_PUBLIC_IPFS_URL}</p>
    </div>
  );
}