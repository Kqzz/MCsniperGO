import React, { useState, useEffect } from "react";
import { Container, Row, Col, Table } from "react-bootstrap";
import ClaimForm from "../components/ClaimForm";
import {
  GetQueues,
  CreateQueue,
  DeleteQueue,
} from "../../wailsjs/go/backendmanager/QueueManager";
import { RemoveButton } from "../components/Buttons";

export default (props) => {
  const [queuedClaims, setQueuedClaims] = useState([]);

  const getQueuedClaims = () => {
    // fetch the queued claims from the backend
    GetQueues().then((res) => {
      setQueuedClaims(res);
    });
  };

  useEffect(getQueuedClaims, []);

  const queueClaim = (claim) => {
    CreateQueue(claim).then(() => {
      getQueuedClaims();
    });
  };

  return (
    <Container>
      <Row>
        <Col>
          <h1>Queue</h1>
        </Col>
      </Row>
      <Row>
        <Col>
          <ClaimForm queueClaim={queueClaim} />
        </Col>
      </Row>
      <Row className="mt-5">
        <Col>
          <h2>Queued Claims</h2>
          <Table responsive>
            <thead>
              <tr>
                <th>Username</th>
                <th>Date</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {queuedClaims.map((claim, index) => (
                <tr key={index}>
                  <td>{claim.username}</td>
                  <td>{claim.startTime}</td>
                  <td>{claim.status}</td>
                  <td>
                    <RemoveButton
                      onClick={() => {
                        DeleteQueue(claim.username);
                        getQueuedClaims();
                      }}
                      data={claim}
                      className="ms-1"
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </Table>
        </Col>
      </Row>
    </Container>
  );
};
