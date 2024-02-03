import {
  Container,
  Flex,
  Text,
  Heading,
  Button,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  useDisclosure,
  Textarea,
  Radio,
  RadioGroup,
  Stack,
} from "@chakra-ui/react";
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

function PlusButton({ onClick }) {
  return (
    <Button onClick={onClick} bgColor={"#1C274C"} width={30} height={30}>
      <Text alt="add" color={"white"}>
        +
      </Text>
    </Button>
  );
}

function AddAccountsModal({ isOpen, onClose, addAccounts }) {
  const [accType, setAccType] = useState("ms");
  return (
    <>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Modal Title</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Textarea
              placeholder={"email:password\nemail:password\netc..."}
              mb={2}
            />

            <RadioGroup onChange={setAccType} value={accType}>
              <Stack direction="row">
                <Radio value="ms" _selected={true}>
                  Microsoft (has username)
                </Radio>
                <Radio value="gc">Giftcard / prename</Radio>
              </Stack>
            </RadioGroup>
          </ModalBody>

          <ModalFooter>
            <Button variant="ghost" onClick={onClose}>
              Close
            </Button>
            <Button colorScheme="blue" mr={2} onClick={onClose}>
              Add Accounts
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}

export default (props) => {
  let [accounts, setAccounts] = useState([]);
  const { isOpen, onOpen, onClose } = useDisclosure();

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

  const addAccounts = () => {
    console.log("add account");
    onOpen();
  };

  return (
    <>
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
          <Flex
            direction={"row"}
            alignItems={"center"}
            justifyContent={"space-between"}
          >
            <Heading>Accounts</Heading>
            <Flex direction={"row"} alignItems={"center"}>
              <Refresh onClick={refreshAccounts} />
              <PlusButton onClick={addAccounts} />
            </Flex>
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
      <AddAccountsModal isOpen={isOpen} onClose={onClose} />
    </>
  );
};
