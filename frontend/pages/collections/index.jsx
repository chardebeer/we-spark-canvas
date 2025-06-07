// pages/collections/index.jsx
import { useState, useEffect } from 'react'
import NextLink from 'next/link'
import Head from 'next/head'
import {
  Box,
  Container,
  Heading,
  Text,
  Button,
  SimpleGrid,
  Flex,
  Link,
  Card,
  CardBody,
  CardHeader,
  CardFooter,
  Skeleton,
} from '@chakra-ui/react'
import { Plus } from 'lucide-react'
import apiClient from '../../lib/axios'

export default function Collections() {
  const [collections, setCollections] = useState([])
  const [loading, setLoading] = useState(true)

  // static fallbacks for bg & border
  const bgColor = 'white'
  const borderColor = 'gray.200'

  useEffect(() => {
    apiClient
      .get('/collections')
      .then((res) => {
        setCollections(Array.isArray(res.data) ? res.data : [])
      })
      .catch((err) => {
        console.error('Error fetching collections:', err)
        alert(err.response?.data?.error || 'Could not load collections')
        setCollections([])
      })
      .finally(() => setLoading(false))
  }, [])

  return (
    <>
      <Head>
        <title>Collections | We Spark Canvas</title>
      </Head>

      <Container maxW="container.xl" py={8}>
        <Flex justify="space-between" align="center" mb={8}>
          <Box>
            <Heading as="h1" size="xl" color="pink.500" mb={2}>
              Collections
            </Heading>
            <Text color="gray.600">
              Discover curated visual inspirations or create your own
            </Text>
          </Box>
          <NextLink href="/collections/new" passHref legacyBehavior>
            <Button
              as="a"
              leftIcon={<Plus size={16} />}
              colorScheme="pink"
              size="md"
            >
              Create Collection
            </Button>
          </NextLink>
        </Flex>

        {loading ? (
          <SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={6}>
            {[...Array(6)].map((_, i) => (
              <Skeleton key={i} height="200px" borderRadius="lg" />
            ))}
          </SimpleGrid>
        ) : collections.length === 0 ? (
          <Box
            textAlign="center"
            p={10}
            borderWidth={1}
            borderRadius="lg"
            borderStyle="dashed"
            borderColor={borderColor}
          >
            <Heading as="h3" size="md" mb={3} color="gray.500">
              No Collections Yet
            </Heading>
            <Text color="gray.500" mb={6}>
              Be the first to create an inspiring collection!
            </Text>
            <NextLink href="/collections/new" passHref legacyBehavior>
              <Button as="a" colorScheme="pink" size="md">
                Create Collection
              </Button>
            </NextLink>
          </Box>
        ) : (
          <SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={6}>
            {collections.map((collection) => (
              <Card
                key={collection.id}
                borderRadius="lg"
                overflow="hidden"
                boxShadow="md"
                bg={bgColor}
                transition="transform 0.3s, box-shadow 0.3s"
                _hover={{ transform: 'translateY(-4px)', boxShadow: 'lg' }}
              >
                <CardHeader pb={0}>
                  <Heading as="h3" size="md" color="pink.500">
                    <NextLink
                      href={`/collections/${collection.id}`}
                      passHref
                      legacyBehavior
                    >
                      <Link _hover={{ textDecoration: 'none' }}>
                        {collection.title}
                      </Link>
                    </NextLink>
                  </Heading>
                </CardHeader>
                <CardBody>
                  <Text noOfLines={2} color="gray.600">
                    {collection.description || 'No description provided'}
                  </Text>
                </CardBody>
                <CardFooter pt={0}>
                  <NextLink
                    href={`/collections/${collection.id}`}
                    passHref
                    legacyBehavior
                  >
                    <Button as="a" variant="outline" colorScheme="pink" size="sm">
                      View Collection
                    </Button>
                  </NextLink>
                </CardFooter>
              </Card>
            ))}
          </SimpleGrid>
        )}
      </Container>
    </>
  )
}
