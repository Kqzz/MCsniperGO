import { Container, Flex, Heading } from "@chakra-ui/react";
import ClaimForm from "../components/ClaimForm";

export default (props) => {
  return (
    <Flex
      as="main"
      role="main"
      direction="column"
      flex="1"
      py="16"
      height="100vh"
      {...props}
      bg=""
      ml={{ base: "0" }}
    >
      <Container maxW={"90%"}>
        <Heading>Queue</Heading>
        <ClaimForm />
      </Container>
    </Flex>
  );
};
