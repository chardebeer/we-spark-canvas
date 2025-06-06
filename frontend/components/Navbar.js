import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import {
  Box,
  Flex,
  Text,
  Button,
  Stack,
  Container,
} from "@chakra-ui/react";

export default function Navbar() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const router = useRouter();
  
  useEffect(() => {
    // Check authentication status on component mount
    const token = localStorage.getItem("token");
    if (token) {
      setIsLoggedIn(true);
    }
  }, []);
  
  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    setIsLoggedIn(false);
    router.push("/");
  };

  return (
    <Box
      bg="white"
      color="gray.600"
      borderBottom="1px"
      borderStyle="solid"
      borderColor="gray.200"
      shadow="sm"
      position="sticky"
      top={0}
      zIndex={10}
    >
      <Container maxW="container.xl">
        <Flex
          minH="60px"
          py={2}
          px={4}
          align="center"
          justify="space-between"
        >
          <Flex flex={1} justify="start">
            <Link href="/" passHref>
              <Text
                cursor="pointer"
                fontFamily="heading"
                fontWeight="bold"
                fontSize="xl"
                color="pink.500"
              >
                We Spark Canvas
              </Text>
            </Link>

            <Flex display={{ base: "none", md: "flex" }} ml={10}>
              <Stack direction="row" spacing={4}>
                <Link href="/" passHref>
                  <Text
                    cursor="pointer"
                    p={2}
                    fontSize="sm"
                    fontWeight={500}
                    color="gray.600"
                    _hover={{ color: "pink.500" }}
                  >
                    Explore
                  </Text>
                </Link>
                <Link href="/collections" passHref>
                  <Text
                    cursor="pointer"
                    p={2}
                    fontSize="sm"
                    fontWeight={500}
                    color="gray.600"
                    _hover={{ color: "pink.500" }}
                  >
                    Collections
                  </Text>
                </Link>
                <Link href="/upload" passHref>
                  <Text
                    cursor="pointer"
                    p={2}
                    fontSize="sm"
                    fontWeight={500}
                    color="gray.600"
                    _hover={{ color: "pink.500" }}
                  >
                    Upload
                  </Text>
                </Link>
              </Stack>
            </Flex>
          </Flex>

          <Stack
            flex={{ base: 1, md: 0 }}
            justify="flex-end"
            direction="row"
            spacing={6}
          >
            {isLoggedIn ? (
              <Button
                onClick={handleLogout}
                fontSize="sm"
                fontWeight={400}
                variant="link"
                color="pink.500"
              >
                Logout
              </Button>
            ) : (
              <>
                <Link href="/login" passHref>
                  <Button
                    as="a"
                    fontSize="sm"
                    fontWeight={400}
                    variant="link"
                    color="pink.500"
                  >
                    Sign In
                  </Button>
                </Link>
                <Link href="/login" passHref>
                  <Button
                    as="a"
                    display={{ base: "none", md: "inline-flex" }}
                    fontSize="sm"
                    fontWeight={600}
                    color="white"
                    bg="pink.400"
                    _hover={{
                      bg: "pink.500",
                    }}
                  >
                    Sign Up
                  </Button>
                </Link>
              </>
            )}
          </Stack>
        </Flex>
      </Container>
    </Box>
  );
}