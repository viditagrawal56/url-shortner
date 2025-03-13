import React, { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Link2, ArrowLeft } from "lucide-react";

const NotFound: React.FC = () => {
  const navigate = useNavigate();
  useEffect(() => {
    // Disable scrolling
    document.body.style.overflow = "hidden";

    return () => {
      // Re-enable scrolling when the component unmounts
      document.body.style.overflow = "";
    };
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-b from-indigo-50 via-white to-purple-50 flex items-center justify-center px-4">
      <div className="max-w-lg w-full text-center">
        <div className="relative mb-8">
          <div className="relative bg-white rounded-xl shadow-xl p-8">
            {/* 404 Number */}
            <h1 className="text-8xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
              404
            </h1>

            {/* Animated Elements */}
            <div className="relative h-40">
              {/* Broken Link Animation */}
              <div className="absolute left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-1/2 flex items-center">
                <Link2 className="h-12 w-12 text-indigo-400 animate-pulse" />
                <div className="w-8 h-2 bg-red-400 rounded-full mx-1 animate-bounce"></div>
                <Link2 className="h-12 w-12 text-purple-400 animate-pulse" />
              </div>
            </div>

            <h2 className="text-2xl font-bold text-gray-900 mb-3">
              Oops! This link seems to be lost
            </h2>
            <p className="text-gray-600 mb-6">
              Even our URL shortener couldn't make this page any shorter! The
              link you're looking for might have been moved, deleted, or never
              existed.
            </p>

            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <button
                onClick={() => navigate("/")}
                className="px-6 py-3 bg-indigo-600 text-white font-semibold rounded-lg shadow-lg hover:bg-indigo-700 transition duration-300 flex items-center justify-center cursor-pointer"
              >
                <ArrowLeft className="mr-2 h-5 w-5" />
                Go Home
              </button>
              <button
                onClick={() => navigate(-1)}
                className="px-6 py-3 bg-white text-indigo-600 font-semibold rounded-lg shadow-md border border-indigo-100 hover:bg-indigo-50 transition duration-300 cursor-pointer"
              >
                Go Back
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NotFound;
