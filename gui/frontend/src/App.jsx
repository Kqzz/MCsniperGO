import logo from "./assets/images/logo-universal.png";
import { Greet } from "../wailsjs/go/main/App";
import { Flex, Box } from "@chakra-ui/react";
import Footer from "./components/Footer";
import Main from "./routes/Main";
import Navbar from "./components/Navbar";

function App() {
  return (
    <Flex direction="column" minH="10vh">
      <Box flex="1">
        <Flex direction={{ base: "row" }}>
          <Box width={{ base: "25%" }}>
            <Navbar />
          </Box>
          <Box width={{ base: "75%" }}>
            <Main />
          </Box>
        </Flex>
      </Box>
      <Footer />
    </Flex>
  );
}

export default App;
