import { PropsWithChildren, ButtonHTMLAttributes } from "react";

const ActionButton = ({
  children,
  ...props
}: PropsWithChildren<ButtonHTMLAttributes<HTMLButtonElement>>) => (
  <button
    className="text-gray-600 hover:text-indigo-600 transition-colors p-1 rounded-md hover:bg-indigo-50 cursor-pointer"
    {...props}
  >
    {children}
  </button>
);

export default ActionButton;
