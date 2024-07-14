import { HashRouter, Routes, Route } from "react-router-dom";
import MainLayout from "./components/Layouts/MainLayout";

import Main from "./routes/Main";
import Accounts from "./routes/Accounts";
import Proxies from "./routes/Proxies";
import Queue from "./routes/Queue";
import Logs from "./routes/Logs";

import "bootstrap/dist/css/bootstrap.min.css";

export default function App() {
  return (
    <HashRouter basename="/">
      <Routes>
        <Route index element={<Main />} />
        <Route path="accounts" element={<Accounts />} />
        <Route path="queue" element={<Queue />} />
        <Route path="proxies" element={<Proxies />} />
        <Route path="logs" element={<Logs />} />
      </Routes>
    </HashRouter>
  );
}
