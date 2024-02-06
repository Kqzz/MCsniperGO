import {
  Flex,
  Text,
  Button
} from "@chakra-ui/react";

import RefreshIcon from "../assets/images/refresh.svg";

function RefreshButton({ onClick, ...params }) {
  return (
    <Flex onClick={onClick} bgColor={"#1C274C"} width={30} height={30} _hover={{ cursor: "pointer" }} style={{
      borderRadius: 5,
        justifyContent: "center",
        alignItems: "center",
        padding: 5,
    }} {...params} >
      <img src={RefreshIcon} alt="refresh" color="white"/>
    </Flex>
  );
}

function RemoveButton({ onClick, data, ...params}) {
  return (
    <Button onClick={() => onClick(data)} color={"red.300"} {...params}>
      <Text>RM</Text>
    </Button>
  );
}

function PlusButton({ onClick, ...params }) {
  return (
    <Flex onClick={onClick} bgColor={"#1C274C"} width={30} height={30} _hover={{ cursor: "pointer" }} style={{
      borderRadius: 5,
        justifyContent: "center",
        alignItems: "center",
        padding: 5,
    }} {...params} >
      <Text alt="add" color={"white"} fontSize="xl">+</Text>
    </Flex>
  );
}

export {RefreshButton, RemoveButton, PlusButton}
