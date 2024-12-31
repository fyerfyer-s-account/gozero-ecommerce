import React from 'react';

interface ButtonProps {
    onClick?: () => void;
    children: React.ReactNode;
    className?: string;
    disabled?: boolean;
}

const Button: React.FC<ButtonProps> = ({ onClick, children, className = '', disabled = false }) => {
    return (
        <button
            onClick={onClick}
            className={`bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 ${className}`}
            disabled={disabled}
        >
            {children}
        </button>
    );
};

export default Button;