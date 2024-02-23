import {
  Container,
  Flex,
  Heading,
  TableContainer,
  Table,
  Thead,
  Th,
  Tr,
  Td,
  Tbody,
} from "@chakra-ui/react";
import ClaimForm from "../components/ClaimForm";
import {
  GetQueues,
  CreateQueue,
  DeleteQueue,
} from "../../wailsjs/go/backendmanager/QueueManager";
import { useState, useEffect } from "react";
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
        <ClaimForm queueClaim={queueClaim} />
        <TableContainer>
          <Heading mt={5}>Queued Claims</Heading>
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Username</Th>
                <Th>Date</Th>
                <Th>Status</Th>
                <Th></Th>
              </Tr>
            </Thead>
            <Tbody>
              {
                // map through the queued claims and render them
                queuedClaims.map((claim, index) => {
                  return (
                    <Tr key={index}>
                      <Td>{claim.username}</Td>
                      <Td>{claim.startTime}</Td>
                      <Td>{claim.status}</Td>
                      <Td>
                        <RemoveButton
                          onClick={() => {
                            DeleteQueue(claim.username);
                            getQueuedClaims();
                          }}
                          data={claim}
                          ml={1}
                        />
                      </Td>
                    </Tr>
                  );
                })
              }
            </Tbody>
          </Table>
        </TableContainer>
      </Container>
    </Flex>
  );
};
