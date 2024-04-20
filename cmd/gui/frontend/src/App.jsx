import { HashRouter, Routes, Route } from "react-router-dom";
import { ChakraProvider, extendBaseTheme } from "@chakra-ui/react";
import MainLayout from "./components/Layouts/MainLayout";

import Main from "./routes/Main";
import Accounts from "./routes/Accounts";
import Proxies from "./routes/Proxies";
import Queue from "./routes/Queue";
import { ChakraUIProvider } from "./chakra-ui/custom-provider";

const config = {
  initialColorMode: "dark",
  useSystemColorMode: false,
};

const theme = extendBaseTheme({ config });

export default function App() {
  return (
    <ChakraUIProvider>
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
    </ChakraUIProvider>
  );
}
