import { Flex, Box } from "@chakra-ui/react";
import Footer from "../Footer";
import Navbar from "../Navbar";

export default ({ children }) => {
  return (
    <Flex direction="column" minH="10vh">
      <Box flex="1">
        <Flex direction={{ base: "row" }}>
          <Box width={{ base: "20%", md: "25%" }}>
            <Navbar />
          </Box>
          <Box width={{ base: "80%", md: "75%" }}>
            <Flex
              as="main"
              role="main"
              direction="column"
              flex="1"
              py="16"
              height="100vh"
              ml={{ base: "0" }}
            >
              {children}
            </Flex>
            <Footer />
          </Box>
        </Flex>
      </Box>
    </Flex>
  );
};
