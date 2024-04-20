import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { ColorModeScript } from "@chakra-ui/react";
import { chakraCustomTheme } from "./chakra-ui/custom-theme";
import App from "./App";

const container = document.getElementById("root");
const root = createRoot(container);

root.render(
  <StrictMode>
    <ColorModeScript
      initialColorMode={chakraCustomTheme.config.initialColorMode}
    />
    <App />
  </StrictMode>
);
