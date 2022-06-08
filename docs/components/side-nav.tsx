import Link from 'next/link';
import { useRouter } from 'next/router';
import { Table } from '../table';

const SideNav: React.FC<{ table: Table }> = ({ table }) => {
    const router = useRouter();

    return (
        <nav className="w-80 flex-shrink-0">
            <ul>
                {table.map(section => (
                    <li key={section.header} className="mb-3">
                        <h5 className="mb-2 font-semibold">{section.header}</h5>
                        <ul>
                            {section.items.map(item => (
                                <li
                                    key={item.name}
                                    className={`border-l py-[3px] pl-4 text-gray-600  ${
                                        router.pathname === item.slug
                                            ? 'border-blue-500 font-medium text-blue-500'
                                            : 'border-gray-300 text-gray-900 hover:border-blue-500'
                                    }`}
                                >
                                    {item.slug ? (
                                        <Link href={item.slug}>
                                            <a className="block">{item.name}</a>
                                        </Link>
                                    ) : (
                                        <a href={item.href!} target="_blank" rel="noopener noreferrer">
                                            {item.name}
                                        </a>
                                    )}
                                </li>
                            ))}
                        </ul>
                    </li>
                ))}
            </ul>
        </nav>
    );
};

export default SideNav;
