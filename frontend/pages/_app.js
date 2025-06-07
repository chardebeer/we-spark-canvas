// pages/_app.js
import Navbar from "../components/Navbar";
import {
  ChakraProvider,
  createSystem,
  defaultConfig,
  defineConfig,
} from "@chakra-ui/react"

const system = createSystem(defaultConfig, {
  theme: { /* your token overrides */ },
})

export default function App({ Component, pageProps }) {
  return (
    <ChakraProvider value={system}>
            <Navbar />

      <Component  {...pageProps} />
    </ChakraProvider>
  )
}
