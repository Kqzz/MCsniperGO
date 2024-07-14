import Footer from "../Footer";
import Navbar from "../Navbar";
import { Container, Row, Col } from "react-bootstrap";

export default ({ children }) => {
  return (
    <Container fluid>
      <Row>
        <Col xs={12} md={3}>
          <Navbar />
        </Col>
        <Col xs={12} md={9}>
          {children}
        </Col>
      </Row>
    </Container>
  );
};
