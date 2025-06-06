import { useState } from 'react';
import { Box, Button, FormControl, FormLabel, Input, Textarea, Stack, Heading } from '@chakra-ui/react';
import { useRouter } from 'next/router';
import axios from '../lib/axios';

export default function CollectionForm() {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      // Get token from localStorage
      const token = localStorage.getItem('token');
      
      if (!token) {
        console.error('Authentication required: Please login to create a collection');
        router.push('/login');
        return;
      }

      // Create collection
      const { data } = await axios.post(
        '/collections', 
        { title, description },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      
      // Show success message
      console.log(`Collection created! Your collection "${title}" was created successfully.`);
      
      // Redirect to the new collection page
      router.push(`/collections/${data.id}`);
    } catch (error) {
      const message = error.response?.data?.error || 'An error occurred';
      console.error('Error creating collection:', message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Box maxW="md" mx="auto" p={6} borderWidth={1} borderRadius="lg" boxShadow="lg">
      <Heading mb={6} textAlign="center" color="pink.400">
        Create New Collection
      </Heading>
      
      <form onSubmit={handleSubmit}>
        <Stack spacing={4}>
          <FormControl id="title" isRequired>
            <FormLabel>Title</FormLabel>
            <Input 
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="My Inspiration Board"
              bg="white"
            />
          </FormControl>
          
          <FormControl id="description">
            <FormLabel>Description</FormLabel>
            <Textarea 
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="A collection of images that inspire me..."
              bg="white"
              resize="vertical"
              rows={4}
            />
          </FormControl>
          
          <Button
            type="submit"
            colorScheme="pink"
            isLoading={isLoading}
            loadingText="Creating collection..."
            w="full"
            mt={4}
          >
            Create Collection
          </Button>
        </Stack>
      </form>
    </Box>
  );
}