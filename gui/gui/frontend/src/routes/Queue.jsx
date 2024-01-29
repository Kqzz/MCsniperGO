import { Container, Flex } from "@chakra-ui/react";

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
      <Container flex="1">Queue</Container>
    </Flex>
  );
};
