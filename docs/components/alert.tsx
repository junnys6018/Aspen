import { FaInfoCircle, FaTimesCircle, FaCheckCircle, FaExclamationTriangle } from 'react-icons/fa';

const Alert: React.FC<{
    children: React.ReactNode;
    level: 'info' | 'success' | 'warning' | 'error';
}> = ({ children, level }) => {
    const icons = {
        info: <FaInfoCircle className="mt-[0.3rem] flex-shrink-0" />,
        success: <FaCheckCircle className="mt-[0.3rem] flex-shrink-0 text-green-600" />,
        warning: <FaExclamationTriangle className="mt-[0.3rem] flex-shrink-0 text-amber-600" />,
        error: <FaTimesCircle className="mt-[0.3rem] flex-shrink-0 text-red-600" />,
    };

    const classNames = {
        info: 'bg-white border',
        success: 'bg-green-100',
        warning: 'bg-amber-100',
        error: 'bg-red-100',
    };

    return (
        <div className={`flew-row not-prose my-5 flex gap-2 rounded-md p-3 ${classNames[level]}`}>
            {icons[level]}
            <div>{children}</div>
        </div>
    );
};

export default Alert;
