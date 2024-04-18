import { Container, Flex, Heading } from "@chakra-ui/react";

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
      <Container flex="1" ml={"5rem"}>
        <Heading>Home</Heading>
        <p>
          Welcome to MCsniperGO! Join the{" "}
          <a href="https://discord.gg/mcsnipergo-734794891258757160">Discord</a>{" "}
          for guides and assistance with using the claimer.
        </p>
      </Container>
    </Flex>
  );
};
