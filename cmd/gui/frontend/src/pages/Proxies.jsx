import {
  AddProxies,
  GetProxies,
  RemoveProxies,
} from "../../wailsjs/go/backendmanager/ProxyManager";
import { useEffect, useState } from "react";
import { PlusButton, RefreshButton, RemoveButton } from "../components/Buttons";
import {
  Container,
  Row,
  Col,
  Table,
  Modal,
  Form,
  Button,
} from "react-bootstrap";

function AddProxiesModal({ show, onHide, addProxies }) {
  const [proxyType, setProxyType] = useState("http");
  const [proxies, setProxies] = useState("");

  return (
    <Modal show={show} onHide={onHide}>
      <Modal.Header closeButton>
        <Modal.Title>Add Proxies</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form.Control
          as="textarea"
          rows={3}
          onChange={(e) => setProxies(e.target.value)}
          placeholder={"10.10.10.10:8080\n10.10.10.11:8080\netc..."}
          className="mb-2"
        />
        <Form.Group>
          <Form.Check
            inline
            type="radio"
            label="HTTP"
            name="proxyType"
            id="http"
            checked={proxyType === "http"}
            onChange={() => setProxyType("http")}
          />
          <Form.Check
            inline
            type="radio"
            label="SOCKS5"
            name="proxyType"
            id="socks5"
            checked={proxyType === "socks5"}
            onChange={() => setProxyType("socks5")}
          />
          <Form.Check
            inline
            type="radio"
            label="SOCKS4"
            name="proxyType"
            id="socks4"
            checked={proxyType === "socks4"}
            onChange={() => setProxyType("socks4")}
          />
        </Form.Group>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={onHide}>
          Close
        </Button>
        <Button
          variant="primary"
          onClick={() => {
            addProxies(proxies, proxyType);
            onHide();
          }}
        >
          Add Proxies
        </Button>
      </Modal.Footer>
    </Modal>
  );
}

export default (props) => {
  let [proxies, setProxies] = useState([]);
  const [showModal, setShowModal] = useState(false);

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
      <Container fluid style={{ maxWidth: "80%", marginLeft: "2rem" }}>
        <Row className="align-items-center justify-content-between mb-3">
          <Col>
            <h1 className="ps-3">Proxies</h1>
          </Col>
          <Col xs="auto">
            <RefreshButton onClick={refreshProxies} className="me-1" />
            <PlusButton onClick={() => setShowModal(true)} />
          </Col>
        </Row>
        <Table responsive>
          <thead>
            <tr>
              <th>Proxy</th>
              <th>Type</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {proxies.map((proxy, index) => (
              <tr key={index}>
                <td>{proxy.url}</td>
                <td>{proxy.type || "N/A"}</td>
                <td>
                  <RemoveButton onClick={removeProxies} data={proxy.url} />
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      </Container>
      <AddProxiesModal
        show={showModal}
        onHide={() => setShowModal(false)}
        addProxies={addProxies}
      />
    </>
  );
};
