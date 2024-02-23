import { Box, Container, Text } from "@chakra-ui/react";

export default (props) => {
  const year = new Date().getFullYear();

  return (
    <Box as="footer" role="contentinfo" bg="bg.accent.default" {...props}>
      <Container>
        <Text style={{ fontSize: "xs" }}>Copyright © {year} Kqzz</Text>
      </Container>
    </Box>
  );
};
