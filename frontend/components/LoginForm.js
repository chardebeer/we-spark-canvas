import { useState } from 'react';
import { Box, Button, FormControl, FormLabel, Input, Stack, Heading, Text, Link } from '@chakra-ui/react';
import { useRouter } from 'next/router';
import axios from '../lib/axios';

export default function LoginForm() {
  const [isLogin, setIsLogin] = useState(true);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [avatarUrl, setAvatarUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const endpoint = isLogin ? '/auth/login' : '/auth/register';
      const payload = isLogin 
        ? { username, password }
        : { username, password, avatar_url: avatarUrl };
      
      const { data } = await axios.post(endpoint, payload);
      
      // Store token in localStorage
      localStorage.setItem('token', data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
      
      // Show success message
      console.log(isLogin ? 'Logged in successfully!' : 'Account created!');
      
      // Redirect to home page
      router.push('/');
    } catch (error) {
      const message = error.response?.data?.error || 'An error occurred';
      console.error('Error:', message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Box maxW="md" mx="auto" p={6} borderWidth={1} borderRadius="lg" boxShadow="lg">
      <Heading mb={6} textAlign="center" color="pink.400">
        {isLogin ? 'Login' : 'Create Account'}
      </Heading>
      
      <form onSubmit={handleSubmit}>
        <Stack spacing={4}>
          <FormControl id="username" isRequired>
            <FormLabel>Username</FormLabel>
            <Input 
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Your username"
              bg="white"
            />
          </FormControl>
          
          <FormControl id="password" isRequired>
            <FormLabel>Password</FormLabel>
            <Input 
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Your password"
              bg="white"
            />
          </FormControl>
          
          {!isLogin && (
            <FormControl id="avatarUrl">
              <FormLabel>Avatar URL (optional)</FormLabel>
              <Input 
                value={avatarUrl}
                onChange={(e) => setAvatarUrl(e.target.value)}
                placeholder="https://example.com/avatar.jpg"
                bg="white"
              />
            </FormControl>
          )}
          
          <Button
            type="submit"
            colorScheme="pink"
            isLoading={isLoading}
            loadingText={isLogin ? "Logging in..." : "Creating account..."}
            w="full"
            mt={4}
          >
            {isLogin ? 'Login' : 'Create Account'}
          </Button>
        </Stack>
      </form>
      
      <Text mt={4} textAlign="center">
        {isLogin ? "Don't have an account?" : "Already have an account?"}
        <Link
          color="pink.500"
          onClick={() => setIsLogin(!isLogin)}
          ml={2}
          cursor="pointer"
        >
          {isLogin ? 'Sign up' : 'Login'}
        </Link>
      </Text>
    </Box>
  );
}