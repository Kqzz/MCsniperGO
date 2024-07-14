import React, { useState } from "react";
import { Container, Form, Row, Col, Button } from "react-bootstrap";

function ClaimForm({ queueClaim }) {
  const [username, setUsername] = useState("");
  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [useProxies, setUseProxies] = useState(true);
  const [runInfinitely, setRunInfinitely] = useState(false);

  const handleSubmit = (event) => {
    event.preventDefault();

    queueClaim({
      username,
      startTime: new Number(startTime),
      endTime: new Number(endTime),
      infinite: runInfinitely,
      useProxies,
    });
  };

  const handleStartTimeChange = (e) => {
    setStartTime(new Number(e.target.value));
  };

  const handleEndTimeChange = (e) => {
    setEndTime(new Number(e.target.value));
  };

  const formatTimestamp = (timestamp) => {
    const date = new Date(timestamp * 1000);
    return date.toLocaleString(); // Adjust the format as per your requirements
  };

  return (
    <Container className="border rounded p-4 mt-3">
      <Form onSubmit={handleSubmit}>
        <Row className="mb-3">
          <Form.Group as={Col} controlId="formUsername">
            <Form.Label>Username</Form.Label>
            <Form.Control
              required
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </Form.Group>
          <Form.Group as={Col} controlId="formStartTime">
            <Form.Label>Start Time (Unix Timestamp)</Form.Label>
            <Form.Control
              required
              type="text"
              value={startTime}
              onChange={handleStartTimeChange}
              disabled={runInfinitely}
            />
          </Form.Group>
          <Form.Group as={Col} controlId="formEndTime">
            <Form.Label>End Time (Unix Timestamp)</Form.Label>
            <Form.Control
              required
              type="text"
              value={endTime}
              onChange={handleEndTimeChange}
              disabled={runInfinitely}
            />
          </Form.Group>
        </Row>
        <Row className="mb-3">
          <Form.Group as={Col} controlId="formFormattedStartTime">
            <Form.Label>Start Time</Form.Label>
            <Form.Control
              type="text"
              value={formatTimestamp(startTime)}
              readOnly
            />
          </Form.Group>
          <Form.Group as={Col} controlId="formFormattedEndTime">
            <Form.Label>End Time</Form.Label>
            <Form.Control
              type="text"
              value={formatTimestamp(endTime)}
              readOnly
            />
          </Form.Group>
        </Row>
        <Row className="mb-3">
          <Col>
            <Form.Check
              type="checkbox"
              label="Use proxies"
              checked={useProxies}
              onChange={(e) => setUseProxies(e.target.checked)}
            />
          </Col>
          <Col>
            <Form.Check
              type="checkbox"
              label="Run infinitely"
              checked={runInfinitely}
              onChange={(e) => setRunInfinitely(e.target.checked)}
            />
          </Col>
          <Col>
            <Button type="submit" variant="primary">
              Queue Claim
            </Button>
          </Col>
        </Row>
      </Form>
    </Container>
  );
}

export default ClaimForm;
