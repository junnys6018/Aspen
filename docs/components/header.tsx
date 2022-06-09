import Link from 'next/link';

const Header: React.FC = () => {
    return (
        <header className="border-b">
            <div className="container">
                <Link href="/">
                    <a className="mt-10 mb-3 inline-block text-2xl font-semibold text-blue-600 sm:text-4xl">Aspen</a>
                </Link>
                <h2 className="ml-4 inline-block text-sm sm:text-base">A toy programming language</h2>
            </div>
        </header>
    );
};

export default Header;
