import { Box, Button, Flex } from "@chakra-ui/react";
import { Img } from "@chakra-ui/react";
import Logo from "../assets/images/logo.png";

function NavButton({ children, ...props }) {
  return (
    <Button variant="outline" mt={2} onClick={() => props.onClick()}>
      {children}
    </Button>
  );
}

export default function Navbar() {
  return (
    <Flex
      p={4}
      flexDirection="column"
      bg={"overlay0-dark"}
      height="100vh"
      width={{ base: "25%", md: "25%" }}
      position="fixed"
    >
      <Box mb={10}>
        <Img src={Logo} width="100%" margin={0}></Img>
      </Box>
      <NavButton colorScheme="teal" variant="outline" mr={3}>
        Home
      </NavButton>
      <NavButton colorScheme="teal" variant="outline" mr={3}>
        Accounts
      </NavButton>
      <NavButton colorScheme="teal" variant="outline" mr={3}>
        Queue
      </NavButton>
    </Flex>
  );
}
