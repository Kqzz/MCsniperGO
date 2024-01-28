import { Box, Container } from "@chakra-ui/react";

export default (props) => {
  const year = new Date().getFullYear();

  return (
    <Box as="footer" role="contentinfo" bg="bg.accent.default" {...props}>
      <Container>
        <p>Copyright © {year} Kqzz</p>
      </Container>
    </Box>
  );
};
