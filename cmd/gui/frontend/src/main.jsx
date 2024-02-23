import React from "react";
import { createRoot } from "react-dom/client";
import { HashRouter, Routes, Route } from "react-router-dom";
import { ChakraProvider, theme } from "@chakra-ui/react";
import MainLayout from "./components/Layouts/MainLayout";

import Main from "./routes/Main";
import Accounts from "./routes/Accounts";
import Proxies from "./routes/Proxies";
import Queue from "./routes/Queue";

const container = document.getElementById("root");
const root = createRoot(container);

root.render(
  <React.StrictMode>
    <ChakraProvider theme={theme}>
      <HashRouter basename="/">
        <Routes>
          <Route
            path="/"
            element={
              <MainLayout>
                <Main />
              </MainLayout>
            }
          ></Route>
          <Route
            path="/accounts"
            element={
              <MainLayout>
                <Accounts />
              </MainLayout>
            }
          ></Route>
          <Route
            path="/queue"
            element={
              <MainLayout>
                <Queue />
              </MainLayout>
            }
          ></Route>
          <Route
            path="/proxies"
            element={
              <MainLayout>
                <Proxies />
              </MainLayout>
            }
          ></Route>
        </Routes>
      </HashRouter>
    </ChakraProvider>
  </React.StrictMode>
);
