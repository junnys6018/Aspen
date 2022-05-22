import type { NextPage } from 'next';
import Head from 'next/head';
import TextEditor from '../components/text-editor';

const Home: NextPage = () => {
    return (
        <div>
            <Head>
                <title>Aspen</title>
                <link rel="icon" href="/favicon.ico" />
            </Head>
            <div className="relative" style={{ width: '800px', height: '500px' }}>
                <TextEditor />
            </div>
        </div>
    );
};

export default Home;
