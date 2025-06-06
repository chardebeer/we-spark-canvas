import { Box, Container, Heading } from '@chakra-ui/react';
import CollectionForm from '../../components/CollectionForm';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';

export default function NewCollectionPage() {
  const router = useRouter();
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  
  useEffect(() => {
    // Check if user is logged in
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login?redirect=/collections/new');
    } else {
      setIsAuthenticated(true);
    }
    setIsLoading(false);
  }, [router]);
  
  if (isLoading) {
    return <Box p={10} textAlign="center">Loading...</Box>;
  }
  
  if (!isAuthenticated) {
    return <Box p={10} textAlign="center">Please login to create a collection</Box>;
  }
  
  return (
    <>
      <Head>
        <title>Create Collection | We Spark Canvas</title>
      </Head>
      <Container maxW="container.md" py={10}>
        <Box textAlign="center" mb={10}>
          <Heading as="h1" size="xl" color="pink.500">
            Create a New Collection
          </Heading>
          <Heading as="h2" size="md" fontWeight="normal" mt={2} color="gray.600">
            Organize and curate your favorite inspirations
          </Heading>
        </Box>
        
        <CollectionForm />
      </Container>
    </>
  );
}