import { Box, Image, Text, Flex, HStack } from '@chakra-ui/react';
import Link from 'next/link';

export default function ImageCard({ image }) {
  const bgColor = 'white';
  const textColor = 'gray.700';
  
  // If the URL already uses ipfs.io, you could replace it with your local gateway:
  const displayUrl = image?.url?.startsWith("ipfs://")
    ? `${process.env.NEXT_PUBLIC_IPFS_URL || 'https://ipfs.io'}/ipfs/${image.url.slice(7)}`
    : image?.url?.startsWith("https://ipfs.io/ipfs/")
    ? image.url.replace("https://ipfs.io/ipfs/", `${process.env.NEXT_PUBLIC_IPFS_URL || 'https://ipfs.io'}/ipfs/`)
    : image?.url || 'https://via.placeholder.com/300x200?text=No+Image';

  if (!image) {
    return null;
  }

  return (
    <Box
      borderRadius="lg"
      overflow="hidden"
      boxShadow="md"
      bg={bgColor}
      transition="transform 0.3s, box-shadow 0.3s"
      _hover={{ transform: 'translateY(-4px)', boxShadow: 'lg' }}
      position="relative"
    >
      <Link href={`/images/${image.id}`} passHref>
        <Box cursor="pointer">
          <Image
            src={displayUrl}
            alt={image.caption || "Image"}
            width="100%"
            height="200px"
            objectFit="cover"
            fallbackSrc="https://via.placeholder.com/300x200?text=Loading..."
          />

          <Box p={3}>
            {image.caption && (
              <Text
                color={textColor}
                fontSize="sm"
                fontWeight="medium"
                noOfLines={2}
                mb={2}
              >
                {image.caption}
              </Text>
            )}

            <Flex justify="space-between" align="center">
              <Text color="gray.500" fontSize="xs">
                {image.hearts || 0} hearts
              </Text>
              
              {image.tags && image.tags.length > 0 && (
                <HStack spacing={1} mt={1} flexWrap="wrap">
                  {image.tags.slice(0, 3).map((tag, index) => (
                    <Box
                      key={index}
                      px={2}
                      py={1}
                      bg="pink.100"
                      color="pink.800"
                      borderRadius="full"
                      fontSize="xs"
                    >
                      {tag}
                    </Box>
                  ))}
                  {image.tags.length > 3 && (
                    <Text fontSize="xs" color="gray.500">
                      +{image.tags.length - 3} more
                    </Text>
                  )}
                </HStack>
              )}
            </Flex>
          </Box>
        </Box>
      </Link>
    </Box>
  );
}