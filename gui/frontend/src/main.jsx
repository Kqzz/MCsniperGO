import React from "react";
import { createRoot } from "react-dom/client";
import App from "./App";
import { HashRouter, Routes, Route } from "react-router-dom";
import { ChakraProvider, theme } from "@chakra-ui/react";

const container = document.getElementById("root");
const root = createRoot(container);

root.render(
  <HashRouter basename="/">
    <Routes>
      <Route path="/"></Route>
    </Routes>
    <React.StrictMode>
      <ChakraProvider theme={theme}>
        <App />
      </ChakraProvider>
    </React.StrictMode>
  </HashRouter>
);
