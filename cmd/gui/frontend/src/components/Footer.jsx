import { Row } from "react-bootstrap";

export default (props) => {
  const year = new Date().getFullYear();

  return (
    <Row as="footer" role="contentinfo" bg="bg.accent.default" {...props}>
      {/* <Container> */}
      <p style={{ fontSize: "xs" }}>Copyright © {year} Kqzz</p>
      {/* </Container> */}
    </Row>
  );
};
