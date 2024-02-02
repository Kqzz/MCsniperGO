import { Container, Flex, Text, Heading, Button } from "@chakra-ui/react";
import {
  GetAccounts,
  RemoveAccountByEmail,
} from "../../wailsjs/go/accountmanager/AccountManager";
import { useEffect, useState } from "react";
import RefreshIcon from "../assets/images/refresh.svg";
import {
  Table,
  Thead,
  Tbody,
  Tfoot,
  Tr,
  Th,
  Td,
  TableCaption,
  TableContainer,
} from "@chakra-ui/react";
import { accountmanager } from "../../wailsjs/go/models";

function AccountStatus(status) {
  let color = "green";
  if (status === "Inactive") color = "red.500";

  return (
    <div
      color={color}
      style={{
        width: "10px",
        height: "10px",
        borderRadius: "50%",
        backgroundColor: color,
      }}
    ></div>
  );
}

function Refresh({ onClick }) {
  return (
    <Flex onClick={onClick} _hover={{ cursor: "pointer" }}>
      <img src={RefreshIcon} alt="refresh" width={30} height={30} />
    </Flex>
  );
}

function RemoveButton({ onClick, email }) {
  return (
    <Button onClick={() => onClick(email)} color={"red.300"}>
      <Text>RM</Text>
    </Button>
  );
}

export default (props) => {
  let [accounts, setAccounts] = useState([]);

  // function to reset accounts
  const refreshAccounts = () => {
    GetAccounts().then((res) => {
      setAccounts(res);
    });
  };

  useEffect(() => {
    refreshAccounts();
  });

  const removeAccount = (email) => {
    RemoveAccountByEmail(email).then((res) => {
      console.debug(res);
      refreshAccounts();
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
      <Container flex="1" ml={"5rem"}>
        <Flex direction={"row"}>
          <Heading>Accounts</Heading>
          <Refresh onClick={refreshAccounts} />
        </Flex>
        <TableContainer>
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Email</Th>
                <Th>Type</Th>
                <Th>Status</Th>
                <Th></Th>
              </Tr>
            </Thead>
            <Tbody>
              {accounts.map((account, index) => {
                return (
                  <Tr key={index}>
                    <Td>{account.email}</Td>
                    <Td>{account.type || "N/A"}</Td>
                    <Td>{AccountStatus(account.status)}</Td>
                    <RemoveButton
                      onClick={removeAccount}
                      email={account.email}
                    />
                  </Tr>
                );
              })}
            </Tbody>
          </Table>
        </TableContainer>
      </Container>
    </Flex>
  );
};
