import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import ProtectedRoute from "./components/ProtectedRoutes";
import Dashboard from "./pages/Dashboard";
import { ToastContainer } from "react-toastify";
import CreateURL from "./pages/CreateURL";
import Navbar from "./components/NavBar";

function App() {
  return (
    <>
      <AuthProvider>
        <Router>
          <div className="min-h-screen bg-gray-50">
            <Navbar />
            <main className="container mx-auto px-4">
              <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route
                  path="/dashboard"
                  element={
                    <ProtectedRoute>
                      <Dashboard />
                    </ProtectedRoute>
                  }
                />
                <Route
                  path="/create"
                  element={
                    <ProtectedRoute>
                      <CreateURL />
                    </ProtectedRoute>
                  }
                />
              </Routes>
            </main>
            <ToastContainer
              autoClose={2000}
              pauseOnHover={true}
              position="bottom-right"
            />
          </div>
        </Router>
      </AuthProvider>
    </>
  );
}

export default App;
