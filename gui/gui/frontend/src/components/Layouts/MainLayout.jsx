import { Flex, Box } from "@chakra-ui/react";
import Navbar from "../Navbar";
import Footer from "../Footer";
import { Children } from "react";

export default ({ children }) => {
  return (
    <Flex direction="column" minH="10vh">
      <Box flex="1">
        <Flex direction={{ base: "row" }}>
          <Box width={{ base: "20%", md: "25%" }}>
            <Navbar />
          </Box>
          <Box width={{ base: "80%", md: "75%" }}>
            {children}

            <Footer />
          </Box>
        </Flex>
      </Box>
    </Flex>
  );
};
