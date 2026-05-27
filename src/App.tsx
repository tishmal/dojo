import { Routes, Route, Navigate } from "react-router-dom";
import { UserProvider } from "./store/user";
import Layout from "./components/Layout";
import Tasks from "./pages/Tasks";
import Profile from "./pages/Profile";
import Raids from "./pages/Raids";
import Sensei from "./pages/Sensei";
import DojoPlus from "./pages/DojoPlus";

export default function App() {
  return (
    <UserProvider>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<Navigate to="/tasks" replace />} />
          <Route path="/tasks" element={<Tasks />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/raids" element={<Raids />} />
          <Route path="/sensei" element={<Sensei />} />
          <Route path="/plus" element={<DojoPlus />} />
          <Route path="*" element={<Navigate to="/tasks" replace />} />
        </Route>
      </Routes>
    </UserProvider>
  );
}