import RefreshIcon from "../assets/images/refresh.svg";
import TrashIcon from "../assets/images/trash.svg";
import { Button } from "react-bootstrap";
import Image from "react-bootstrap/Image";

function RefreshButton({ onClick, ...params }) {
  return (
    <Button
      onClick={onClick}
      bgColor={"#1C274C"}
      width={30}
      height={30}
      _hover={{ cursor: "pointer" }}
      style={{
        borderRadius: 5,
        justifyContent: "center",
        alignItems: "center",
        padding: 5,
      }}
      {...params}
    >
      <img src={RefreshIcon} alt="refresh" color="white" />
    </Button>
  );
}

function RemoveButton({ onClick, data, ...params }) {
  return (
    <Button
      onClick={() => onClick(data)}
      color={"red.300"}
      colorScheme="blue"
      marginTop={""}
      {...params}
    >
      <Image src={TrashIcon} alt="remove" width={25} />
    </Button>
  );
}

function PlusButton({ onClick, ...params }) {
  return (
    <Button
      onClick={onClick}
      bgColor={"#1C274C"}
      width={30}
      height={30}
      _hover={{ cursor: "pointer" }}
      style={{
        borderRadius: 5,
        justifyContent: "center",
        alignItems: "center",
        padding: 5,
      }}
      {...params}
    >
      <p alt="add" color={"white"} fontSize="xl">
        +
      </p>
    </Button>
  );
}

export { RefreshButton, RemoveButton, PlusButton };
