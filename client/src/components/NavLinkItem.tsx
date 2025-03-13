import { Link } from "react-router-dom";

interface NavLinkItemProps {
  to: string;
  isActive: boolean;
  icon: React.ElementType;
  label: string;
  onClick?: () => void;
}

function NavLinkItem({
  to,
  isActive,
  icon: Icon,
  label,
  onClick,
}: NavLinkItemProps) {
  return (
    <Link
      to={to}
      onClick={onClick}
      className={`flex items-center px-3 py-2 rounded-md text-sm font-medium transition-colors duration-200 ${
        isActive
          ? "text-indigo-700 bg-indigo-50"
          : "text-gray-700 hover:text-indigo-600 hover:bg-indigo-50"
      }`}
    >
      <Icon className="h-4 w-4 mr-1" />
      {label}
    </Link>
  );
}

export default NavLinkItem;
