import {
  Flex,
  Text,
  Button
} from "@chakra-ui/react";

import RefreshIcon from "../assets/images/refresh.svg";

function RefreshButton({ onClick }) {
  return (
    <Flex onClick={onClick} _hover={{ cursor: "pointer" }}>
      <img src={RefreshIcon} alt="refresh" width={30} height={30} />
    </Flex>
  );
}

function RemoveButton({ onClick, data }) {
  return (
    <Button onClick={() => onClick(data)} color={"red.300"}>
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

export {RefreshButton, RemoveButton, PlusButton}
