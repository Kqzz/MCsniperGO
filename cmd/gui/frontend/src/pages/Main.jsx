import { Container } from "react-bootstrap";

export default (props) => {
  return (
    <Container flex="1" ml={"2rem"}>
      <h1>Home</h1>
      <p>
        Welcome to MCsniperGO! Join the{" "}
        <a href="https://discord.gg/mcsnipergo-734794891258757160">Discord</a>{" "}
        for guides and assistance with using the claimer.
      </p>
    </Container>
  );
};
