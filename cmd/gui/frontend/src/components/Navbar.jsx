import Logo from "../assets/images/logo.png";

import ArticleIcon from "@mui/icons-material/Article";
import QueueIcon from "@mui/icons-material/Queue";
import CloudIcon from "@mui/icons-material/Cloud";
import GroupIcon from "@mui/icons-material/Group";
import HomeIcon from "@mui/icons-material/Home";

import { Link } from "react-router-dom";
import { Container, Button, Image, Row } from "react-bootstrap";

function NavButton({ children, path, leftIcon }) {
  return (
    <Link to={path} mt={2}>
      <Button
        variant="outline"
        color={"gray.50"}
        width={"100%"}
        _hover={{ bg: "gray.700" }}
      >
        {leftIcon}
        {children}
      </Button>
    </Link>
  );
}
export default function Navbar() {
  return (
    <Row>
      <Container mb={10}>
        <Image src={Logo} width="100%" margin={0}></Image>
      </Container>
      <NavButton
        colorScheme="teal"
        variant="outline"
        mr={3}
        path={"/"}
        leftIcon={<HomeIcon />}
      >
        Home
      </NavButton>
      <NavButton
        colorScheme="teal"
        variant="outline"
        mr={3}
        path={"/accounts"}
        leftIcon={<GroupIcon />}
      >
        Accounts
      </NavButton>
      <NavButton
        colorScheme="teal"
        variant="outline"
        mr={3}
        path={"/proxies"}
        leftIcon={<CloudIcon />}
      >
        Proxies
      </NavButton>
      <NavButton
        colorScheme="teal"
        variant="outline"
        mr={3}
        path={"/queue"}
        leftIcon={<QueueIcon />}
      >
        Queue
      </NavButton>
      <NavButton
        colorScheme="teal"
        variant="outline"
        mr={3}
        path={"/logs"}
        leftIcon={<ArticleIcon />}
      >
        Logs
      </NavButton>
    </Row>
  );
}
