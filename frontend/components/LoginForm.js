// components/LoginForm.jsx
import { useState } from "react";
import {
  Box,
  Button,
  Field,
  Input,
  Stack,
  Heading,
  Text,
  Link,
} from "@chakra-ui/react";
import { useRouter } from "next/router";
import axios from "../lib/axios";

export default function LoginForm() {
  const [isLogin, setIsLogin] = useState(true);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [avatarUrl, setAvatarUrl] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const endpoint = isLogin ? "/auth/login" : "/auth/register";
      const payload = isLogin
        ? { username, password }
        : { username, password, avatar_url: avatarUrl };

      const { data } = await axios.post(endpoint, payload);

      localStorage.setItem("token", data.token);
      localStorage.setItem("user", JSON.stringify(data.user));

      router.push("/");
    } catch (error) {
      console.error(
        "Error:",
        error.response?.data?.error || "An error occurred"
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Box maxW="md" mx="auto" p={6} borderWidth={1} borderRadius="lg" boxShadow="lg">
      <Heading mb={6} textAlign="center" color="pink.400">
        {isLogin ? "Login" : "Create Account"}
      </Heading>

      <form onSubmit={handleSubmit}>
        <Stack spacing={4}>
          <Field.Root id="username" required>
            <Field.Label>Username</Field.Label>
            <Input
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Your username"
              bg="white"
            />
          </Field.Root>

          <Field.Root id="password" required>
            <Field.Label>Password</Field.Label>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Your password"
              bg="white"
            />
          </Field.Root>

          {!isLogin && (
            <Field.Root id="avatarUrl">
              <Field.Label>Avatar URL (optional)</Field.Label>
              <Input
                value={avatarUrl}
                onChange={(e) => setAvatarUrl(e.target.value)}
                placeholder="https://example.com/avatar.jpg"
                bg="white"
              />
            </Field.Root>
          )}

          <Button
            type="submit"
            colorScheme="pink"
            isLoading={isLoading}
            loadingText={isLogin ? "Logging in..." : "Creating account..."}
            w="full"
            mt={4}
          >
            {isLogin ? "Login" : "Create Account"}
          </Button>
        </Stack>
      </form>

      <Text mt={4} textAlign="center">
        {isLogin ? "Don't have an account?" : "Already have an account?"}{" "}
        <Link
          color="pink.500"
          onClick={() => setIsLogin(!isLogin)}
          ml={2}
          cursor="pointer"
        >
          {isLogin ? "Sign up" : "Login"}
        </Link>
      </Text>
    </Box>
  );
}
