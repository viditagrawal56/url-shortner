import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import {
  Shield,
  Menu,
  X,
  Home,
  LayoutDashboard,
  Link2,
  LogOut,
  LogIn,
  UserPlus,
} from "lucide-react";
import NavLinkItem from "./NavLinkItem";

export default function Navbar() {
  const { isAuthenticated, logout } = useAuth();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const location = useLocation();

  const toggleMenu = () => setIsMenuOpen((prev) => !prev);
  const closeMenu = () => setIsMenuOpen(false);
  const isActive = (path: string) => location.pathname === path;

  const AuthButtons = () =>
    isAuthenticated ? (
      <>
        <NavLinkItem
          to="/dashboard"
          isActive={isActive("/dashboard")}
          icon={LayoutDashboard}
          label="Dashboard"
          onClick={closeMenu}
        />
        <NavLinkItem
          to="/create"
          isActive={isActive("/create")}
          icon={Link2}
          label="Create URL"
          onClick={closeMenu}
        />
        <button
          onClick={() => {
            logout();
            closeMenu();
          }}
          className="ml-2 px-4 py-2 flex items-center bg-white border border-red-500 text-red-500 rounded-md text-sm font-medium hover:bg-red-50 transition-colors duration-200 cursor-pointer"
        >
          <LogOut className="h-4 w-4 mr-1" />
          Logout
        </button>
      </>
    ) : (
      <>
        <Link
          to="/login"
          className="ml-2 px-4 py-2 flex items-center text-indigo-600 border border-indigo-200 rounded-md text-sm font-medium hover:bg-indigo-50 transition-colors duration-200"
          onClick={closeMenu}
        >
          <LogIn className="h-4 w-4 mr-1" />
          Login
        </Link>
        <Link
          to="/register"
          className="ml-2 px-4 py-2 flex items-center bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-md text-sm font-medium hover:from-indigo-700 hover:to-purple-700 transition-colors duration-200 shadow-sm"
          onClick={closeMenu}
        >
          <UserPlus className="h-4 w-4 mr-1" />
          Register
        </Link>
      </>
    );

  return (
    <nav className="bg-white shadow-md sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          {/* Logo */}
          <div className="flex items-center">
            <Link to="/" className="flex items-center" onClick={closeMenu}>
              <Shield className="h-8 w-8 text-indigo-600 mr-2" />
              <span className="text-xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                SecureLink
              </span>
            </Link>
          </div>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-1">
            <NavLinkItem
              to="/"
              isActive={isActive("/")}
              icon={Home}
              label="Home"
            />
            <AuthButtons />
          </div>

          {/* Mobile Menu Button */}
          <div className="md:hidden flex items-center">
            <button
              onClick={toggleMenu}
              className="inline-flex items-center justify-center p-2 rounded-md text-gray-700 hover:text-indigo-600 hover:bg-indigo-50 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500 cursor-pointer"
            >
              <span className="sr-only">Open main menu</span>
              {isMenuOpen ? (
                <X className="block h-6 w-6" aria-hidden="true" />
              ) : (
                <Menu className="block h-6 w-6" aria-hidden="true" />
              )}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Menu */}
      {isMenuOpen && (
        <div className="md:hidden bg-white shadow-lg rounded-b-lg px-2 pt-2 pb-3 space-y-1 sm:px-3">
          <NavLinkItem
            to="/"
            isActive={isActive("/")}
            icon={Home}
            label="Home"
            onClick={closeMenu}
          />
          <AuthButtons />
        </div>
      )}
    </nav>
  );
}
