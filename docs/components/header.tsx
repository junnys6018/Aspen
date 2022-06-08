import Link from 'next/link';

const Header: React.FC = () => {
    return (
        <header className="border-b">
            <div className="container">
                <Link href="/">
                    <a className="inline-block">
                        <h1 className="mt-10 mb-3 text-4xl font-semibold text-blue-600">Aspen</h1>
                    </a>
                </Link>
                <h2 className="ml-4 inline-block">A toy programming language</h2>
            </div>
        </header>
    );
};

export default Header;
