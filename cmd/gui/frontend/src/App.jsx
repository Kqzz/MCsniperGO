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
      <MainLayout>
        <HashRouter basename="/">
          <Routes>
            <Route path="/" element={<Main />}></Route>
            <Route path="/accounts" element={<Accounts />}></Route>
            <Route path="/queue" element={<Queue />}></Route>
            <Route path="/proxies" element={<Proxies />}></Route>
          </Routes>
        </HashRouter>
      </MainLayout>
    </ChakraUIProvider>
  );
}
