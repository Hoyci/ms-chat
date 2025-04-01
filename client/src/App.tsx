import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { Login } from "@pages/Auth/Login";
import { Signup } from "@pages/Auth/Signup";
import ChatApp from "@pages/ChatApp";
import ProtectedRoute from "@pages/ProtectedRoute";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />

        <Route element={<ProtectedRoute />}>
          <Route path="/" element={<ChatApp />} />
        </Route>

        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
