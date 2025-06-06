import { useState, useEffect } from 'react';
import Head from 'next/head';
import {
  Container,
  Heading,
  SimpleGrid,
  Box,
  Text,
  Skeleton,
  Button,
  Flex
} from '@chakra-ui/react';
import ImageCard from "../components/ImageCard";
import useImages from "../hooks/useImages";
import Link from 'next/link';

export default function Home() {
  const { images, loading, error } = useImages(20, 0);
  const bgColor = "white";

  return (
    <>
      <Head>
        <title>We Spark Canvas - Visual Inspiration</title>
        <meta name="description" content="A beautiful platform for visual inspiration and creativity" />
      </Head>

      <Container maxW="container.xl" py={8}>
        <Flex justify="space-between" align="center" mb={8}>
          <Box>
            <Heading as="h1" size="xl" color="pink.500" mb={2}>
              Visual Inspiration
            </Heading>
            <Text color="gray.600">
              Discover and share beautiful visuals that spark creativity
            </Text>
          </Box>
          <Link href="/upload" passHref>
            <Button
              as="a"
              colorScheme="pink"
              variant="solid"
              size="md"
            >
              Upload Image
            </Button>
          </Link>
        </Flex>

        {loading ? (
          <SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={6}>
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <Skeleton key={i} height="200px" borderRadius="lg" />
            ))}
          </SimpleGrid>
        ) : error ? (
          <Box textAlign="center" p={10}>
            <Heading as="h3" size="md" mb={3} color="red.500">
              Error Loading Images
            </Heading>
            <Text color="gray.500">
              We encountered a problem while loading images. Please try again later.
            </Text>
          </Box>
        ) : images && images.length > 0 ? (
          <SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={6}>
            {images.map((img) => (
              <ImageCard key={img.id} image={img} />
            ))}
          </SimpleGrid>
        ) : (
          <Box 
            textAlign="center" 
            p={10} 
            borderWidth={1} 
            borderRadius="lg" 
            borderStyle="dashed"
          >
            <Heading as="h3" size="md" mb={3} color="gray.500">
              No Images Yet
            </Heading>
            <Text color="gray.500" mb={6}>
              Be the first to upload an inspiring image!
            </Text>
            <Link href="/upload" passHref>
              <Button as="a" colorScheme="pink" size="md">
                Upload Image
              </Button>
            </Link>
          </Box>
        )}
      </Container>
    </>
  );
}
