import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

function Home() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();

  const handleButtonClick = () => {
    if (isAuthenticated) {
      navigate("/dashboard");
    } else {
      navigate("/login");
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-b from-blue-50 to-white">
      <div className="text-center max-w-2xl px-4">
        <h1 className="text-4xl font-bold text-gray-800 mb-6">
          Welcome to URL Shortener
        </h1>
        <p className="text-gray-600 text-lg mb-8">
          Transform your long URLs into short, memorable links in seconds.
        </p>
        <button
          onClick={handleButtonClick}
          className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-8 rounded-lg transition duration-300 cursor-pointer"
        >
          {isAuthenticated ? "Go to Dashboard" : "Get Started"}
        </button>
      </div>
    </div>
  );
}

export default Home;
