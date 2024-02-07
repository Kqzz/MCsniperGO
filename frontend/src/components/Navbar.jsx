import { Box, Button, Flex } from "@chakra-ui/react";
import { Img } from "@chakra-ui/react";
import Logo from "../assets/images/logo.png";

function NavButton({ children, path }) {
  return (
    <a href={path}>
      <Button
        variant="outline"
        mt={2}
        color={"gray.50"}
        width={"100%"}
        _hover={{ bg: "gray.700" }}
      >
        {children}
      </Button>
    </a>
  );
}
export default function Navbar() {
  return (
    <Flex
      p={4}
      flexDirection="column"
      bg={"gray.600"}
      height="100vh"
      width={{ base: "25%" }}
      position="fixed"
    >
      <Box mb={10}>
        <Img src={Logo} width="100%" margin={0}></Img>
      </Box>
      <NavButton colorScheme="teal" variant="outline" mr={3} path={"#"}>
        Home
      </NavButton>
      <NavButton
        colorScheme="teal"
        variant="outline"
        mr={3}
        path={"/#accounts"}
      >
        Accounts
      </NavButton>
      <NavButton colorScheme="teal" variant="outline" mr={3} path={"/#proxies"}>
        Proxies
      </NavButton>
      <NavButton colorScheme="teal" variant="outline" mr={3} path={"/#queue"}>
        Queue
      </NavButton>
    </Flex>
  );
}
