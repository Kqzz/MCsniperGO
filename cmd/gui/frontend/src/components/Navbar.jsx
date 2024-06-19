import { Box, Button, Flex } from "@chakra-ui/react";
import { Img, Link } from "@chakra-ui/react";
import Logo from "../assets/images/logo.png";

import ArticleIcon from "@mui/icons-material/Article";
import QueueIcon from "@mui/icons-material/Queue";
import CloudIcon from "@mui/icons-material/Cloud";
import GroupIcon from "@mui/icons-material/Group";
import HomeIcon from "@mui/icons-material/Home";

function NavButton({ children, path, leftIcon }) {
  return (
    <Link href={path} mt={2}>
      <Button
        variant="outline"
        color={"gray.50"}
        width={"100%"}
        _hover={{ bg: "gray.700" }}
        leftIcon={leftIcon}
      >
        {children}
      </Button>
    </Link>
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
      <NavButton
        colorScheme="teal"
        variant="outline"
        mr={3}
        path={"#"}
        leftIcon={<HomeIcon />}
      >
        Home
      </NavButton>
      <NavButton
        colorScheme="teal"
        variant="outline"
        mr={3}
        path={"/#accounts"}
        leftIcon={<GroupIcon />}
      >
        Accounts
      </NavButton>
      <NavButton
        colorScheme="teal"
        variant="outline"
        mr={3}
        path={"/#proxies"}
        leftIcon={<CloudIcon />}
      >
        Proxies
      </NavButton>
      <NavButton
        colorScheme="teal"
        variant="outline"
        mr={3}
        path={"/#queue"}
        leftIcon={<QueueIcon />}
      >
        Queue
      </NavButton>
      <NavButton
        colorScheme="teal"
        variant="outline"
        mr={3}
        path={"/#logs"}
        leftIcon={<ArticleIcon />}
      >
        Logs
      </NavButton>
    </Flex>
  );
}
