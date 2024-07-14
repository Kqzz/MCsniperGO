import {
  AddAccounts,
  GetAccounts,
  RemoveAccountByEmail,
} from "../../wailsjs/go/backendmanager/AccountManager";

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

function AccountStatus(status) {
  let color = "green";
  if (status === "Inactive") color = "red";

  return (
    <div
      style={{
        width: "10px",
        height: "10px",
        borderRadius: "50%",
        backgroundColor: color,
      }}
    ></div>
  );
}

function AddAccountsModal({ show, onHide, addAccounts }) {
  const [accType, setAccType] = useState("ms");
  const [accounts, setAccounts] = useState("");

  return (
    <Modal show={show} onHide={onHide}>
      <Modal.Header closeButton>
        <Modal.Title>Add Accounts</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form.Control
          as="textarea"
          rows={3}
          onChange={(e) => setAccounts(e.target.value)}
          placeholder={"email:password\nemail:password\netc..."}
          className="mb-2"
        />
        <Form.Check
          type="radio"
          label="Microsoft (has username)"
          name="accType"
          id="ms"
          checked={accType === "ms"}
          onChange={() => setAccType("ms")}
          inline
        />
        <Form.Check
          type="radio"
          label="Giftcard / prename"
          name="accType"
          id="gc"
          checked={accType === "gc"}
          onChange={() => setAccType("gc")}
          inline
        />
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={onHide}>
          Close
        </Button>
        <Button
          variant="primary"
          onClick={() => {
            addAccounts(accounts, accType);
            onHide();
          }}
        >
          Add Accounts
        </Button>
      </Modal.Footer>
    </Modal>
  );
}

export default (props) => {
  let [accounts, setAccounts] = useState([]);
  const [showModal, setShowModal] = useState(false);

  const refreshAccounts = () => {
    GetAccounts().then((res) => {
      setAccounts(res);
    });
  };

  useEffect(() => {
    refreshAccounts();
  }, []);

  const removeAccount = (email) => {
    RemoveAccountByEmail(email).then((res) => {
      console.debug(res);
      refreshAccounts();
    });
  };

  const addAccounts = (accountsString, type) => {
    AddAccounts(accountsString, type).then((res) => {
      console.log(res);
      // TODO notifications
      refreshAccounts();
    });
  };

  return (
    <>
      <Container fluid className="ml-3" style={{ maxWidth: "80%" }}>
        <Row className="align-items-center justify-content-between">
          <Col>
            <h1 className="pl-3">Accounts</h1>
          </Col>
          <Col xs="auto">
            <RefreshButton onClick={refreshAccounts} className="mr-1" />
            <PlusButton onClick={() => setShowModal(true)} />
          </Col>
        </Row>
        <Table responsive>
          <thead>
            <tr>
              <th>Email</th>
              <th>Type</th>
              <th>Status</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {accounts.map((account, index) => (
              <tr key={index}>
                <td>{account.email}</td>
                <td>{account.type || "N/A"}</td>
                <td>{AccountStatus(account.status)}</td>
                <td>
                  <RemoveButton onClick={removeAccount} data={account.email} />
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      </Container>
      <AddAccountsModal
        show={showModal}
        onHide={() => setShowModal(false)}
        addAccounts={addAccounts}
      />
    </>
  );
};
