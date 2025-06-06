import { Box, Container, Heading } from '@chakra-ui/react';
import LoginForm from '../components/LoginForm';
import { useEffect } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';

export default function LoginPage() {
  const router = useRouter();
  
  useEffect(() => {
    // Check if user is already logged in
    const token = localStorage.getItem('token');
    if (token) {
      router.push('/');
    }
  }, [router]);
  
  return (
    <>
      <Head>
        <title>Login | We Spark Canvas</title>
      </Head>
      <Container maxW="container.md" py={10}>
        <Box textAlign="center" mb={10}>
          <Heading as="h1" size="xl" color="pink.500">
            Welcome to We Spark Canvas
          </Heading>
          <Heading as="h2" size="md" fontWeight="normal" mt={2} color="gray.600">
            Login or create an account to start your visual journey
          </Heading>
        </Box>
        
        <LoginForm />
      </Container>
    </>
  );
}