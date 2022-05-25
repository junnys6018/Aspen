import Link from 'next/link';

const Header: React.FC = () => {
    return (
        <header className="border-b">
            <div className="container">
                <h1 className="mt-10 mb-3 text-4xl font-semibold text-blue-600">Aspen</h1>
                <div className="flex flex-row">
                    <h2 className="mr-auto">A toy programming language</h2>
                    <nav className="flex flex-row gap-1">
                        <Link href="/">
                            <a className="text-blue-600 hover:underline">Home</a>
                        </Link>
                        <span className="text-blue-600">|</span>
                        <Link href="/playground">
                            <a className="text-blue-600 hover:underline">Playground</a>
                        </Link>
                        <span className="text-blue-600">|</span>
                        <a
                            className="text-blue-600 hover:underline"
                            href="https://github.com/junnys6018/Aspen"
                            target="_blank"
                            rel="noopener noreferrer"
                        >
                            GitHub
                        </a>
                    </nav>
                </div>
            </div>
        </header>
    );
};

export default Header;
