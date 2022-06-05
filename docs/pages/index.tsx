import type { NextPage } from 'next';
import Link from 'next/link';

const Home: NextPage = () => {
    return (
        <p>
            This page is currently under construction, why dont you checkout our{' '}
            <Link href="/playground">
                <a className="text-blue-600 hover:underline">playgound</a>
            </Link>{' '}
            in the mean time?
        </p>
    );
};

export default Home;
