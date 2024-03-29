import React from "react";

type ButtonProps = {
  className?: string;
  children?: React.ReactNode;
  onClick?: () => void;
};

const Button: React.FC<ButtonProps> = ({ className, children, onClick }) => {
  return (
    <div
      onClick={() => onClick && onClick()}
      className={`cursor-pointer hover:bg-white hover:text-black flex active:opacity-80
            items-center justify-center py-2 px-6 m-2 bg-dark-gray min-w-[8rem]
            w-fit rounded-xl text-base font-bold ${className}`}
    >
      {children}
    </div>
  );
};

export default Button;
