import {
  Container,
  Flex,
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
  AddProxies,
  GetProxies,
  RemoveProxies,
} from "../../wailsjs/go/backendmanager/ProxyManager";

import { useEffect, useState } from "react";
import { PlusButton, RefreshButton, RemoveButton } from "../components/Buttons";
import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
} from "@chakra-ui/react";

function AddProxiesModal({ isOpen, onClose, addProxies }) {
  const [proxyType, setProxyType] = useState("http");
  const [proxies, setProxies] = useState("");

  return (
    <>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Add Accounhts</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Textarea
              onChange={(e) => setProxies(e.target.value)}
              placeholder={"10.10.10.10:8080\n10.10.10.11:8080\netc..."}
              mb={2}
            />

            <RadioGroup onChange={setProxyType} value={proxyType}>
              <Stack direction="row">
                <Radio value="http" _selected={true}>
                  HTTP
                </Radio>
                <Radio value="socks5">SOCKS5</Radio>
                <Radio value="socks4">SOCKS4</Radio>
              </Stack>
            </RadioGroup>
          </ModalBody>

          <ModalFooter>
            <Button variant="ghost" onClick={onClose}>
              Close
            </Button>
            <Button
              colorScheme="blue"
              mr={2}
              onClick={() => {
                addProxies(proxies, proxyType);
                onClose();
              }}
            >
              Add Proxies
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}

export default (props) => {
  let [proxies, setProxies] = useState([]);
  const { isOpen, onOpen, onClose } = useDisclosure();

  // function to reset proxies
  const refreshProxies = () => {
    console.log("refreshing proxies");
    GetProxies().then((res) => {
      if (res === null) return;
      setProxies(res);
    });
  };

  useEffect(() => {
    refreshProxies();
  }, []);

  const removeProxies = (proxies) => {
    // TODO multi proxy removal
    proxies = [proxies];
    RemoveProxies(proxies).then((res) => {
      console.debug(res);
      refreshProxies();
    });
  };

  const addProxiesModalOpen = () => {
    onOpen();
  };

  const addProxies = (proxyString, type) => {
    const proxies = proxyString.split("\n").filter((proxy) => proxy !== "");
    console.log(proxies);
    AddProxies(proxies, type).then((res) => {
      console.log(res);
      // TODO notifications
      refreshProxies();
    });
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
            <Heading pl={"6"}>Proxies</Heading>
            <Flex direction={"row"} alignItems={"center"}>
              <RefreshButton onClick={refreshProxies} mr={1} />
              <PlusButton onClick={addProxiesModalOpen} />
            </Flex>
          </Flex>
          <TableContainer>
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>Proxy</Th>
                  <Th>Type</Th>
                </Tr>
              </Thead>
              <Tbody>
                {proxies.map((proxy, index) => {
                  return (
                    <Tr key={index}>
                      <Td>{proxy.url}</Td>
                      <Td>{proxy.type || "N/A"}</Td>
                      <RemoveButton onClick={removeProxies} data={proxy.url} />
                    </Tr>
                  );
                })}
              </Tbody>
            </Table>
          </TableContainer>
        </Container>
      </Flex>
      <AddProxiesModal
        isOpen={isOpen}
        onClose={onClose}
        addProxies={addProxies}
      />
    </>
  );
};
