import "@/styles/globals.css";
import { ChakraProvider, defaultSystem } from "@chakra-ui/react";
import Navbar from "../components/Navbar";

export default function App({ Component, pageProps }) {
  return (
    <ChakraProvider value={defaultSystem}>
      <Navbar />
      <Component {...pageProps} />
    </ChakraProvider>
  );
}